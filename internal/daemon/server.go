package daemon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/asdine/storm/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	logFilePerm                    = 0o660
	failRestartDelay time.Duration = time.Second * 5
)

var emptyResponse = &emptypb.Empty{}

type DaemonServerOptions struct {
	LogsDirpath string
	DBFactory   func() *storm.DB
}

type DaemonServer struct {
	pb.UnimplementedProcessServiceServer

	Options DaemonServerOptions

	units   map[uint64]*Unit
	unitsMu sync.RWMutex
}

func NewDaemonServer(options DaemonServerOptions) *DaemonServer {
	return &DaemonServer{
		Options: options,
		units:   make(map[uint64]*Unit),
	}
}

func (s *DaemonServer) getUnitLogFilepath(unitID uint64) string {
	return path.Join(s.Options.LogsDirpath, fmt.Sprintf("%d.log", unitID))
}

func (s *DaemonServer) openUnitLogFile(unitID uint64) (*os.File, error) {
	const logFileFlag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	logFilepath := s.getUnitLogFilepath(unitID)
	return os.OpenFile(logFilepath, logFileFlag, logFilePerm)
}

func (s *DaemonServer) watchUnit(unit *Unit) {
	unit.Command.Wait()
	unit.LogFile.Close()

	if unit.Status() != UnitStatusFailed {
		return
	}

	s.unitsMu.RLock()

	if s.units[unit.Model.ID] == nil {
		s.unitsMu.RUnlock()
		return
	}

	s.unitsMu.RUnlock()

	time.Sleep(failRestartDelay)
	s.StartUnit(unit.Model)
}

func (s *DaemonServer) addUnit(unit *Unit) {
	s.unitsMu.Lock()
	s.units[unit.Model.ID] = unit
	s.unitsMu.Unlock()
	go s.watchUnit(unit)
}

func (s *DaemonServer) StartUnit(model UnitModel) (*Unit, error) {
	logFile, err := s.openUnitLogFile(model.ID)

	if err != nil {
		return nil, err
	}

	unitCtx, cancel := context.WithCancel(context.Background())
	command := createUnitStartCommand(unitCtx, &model, logFile)

	if err := command.Start(); err != nil {
		cancel()
		return nil, err
	}

	unit := &Unit{
		Model:     model,
		Command:   command,
		LogFile:   logFile,
		StartedAt: time.Now(),
		Cancel:    cancel,
	}

	s.addUnit(unit)
	return unit, nil
}

func (s *DaemonServer) stopUnitsStream(
	unitIDs []uint64,
	stream pb.ProcessService_StopServer,
) error {
	eg, _ := errgroup.WithContext(stream.Context())

	for _, unitID := range unitIDs {
		id := unitID

		eg.Go(func() error {
			s.unitsMu.RLock()
			unit := s.units[id]
			s.unitsMu.RUnlock()

			response := pb.StopResponse{UnitId: id}

			if unit == nil {
				response.Error = fmt.Sprintf("unit %d not found", id)
				return stream.Send(&response)
			}

			if unit.Status() != UnitStatusRunning {
				response.Error = fmt.Sprintf(
					"unit %s (%d) is not running",
					unit.Model.Name,
					unit.Model.ID,
				)

				response.Unit = unit.PB()
				return stream.Send(&response)
			}

			unit.Stop()
			response.Unit = unit.PB()
			return stream.Send(&response)
		})
	}

	return eg.Wait()
}

func (s *DaemonServer) restartUnitsStream(
	unitIDs []uint64,
	stream pb.ProcessService_RestartServer,
) error {
	db := s.Options.DBFactory()

	defer db.Commit()
	defer db.Close()

	eg, _ := errgroup.WithContext(stream.Context())

	for _, unitID := range unitIDs {
		id := unitID

		eg.Go(func() error {
			s.unitsMu.RLock()
			unit := s.units[id]
			s.unitsMu.RUnlock()

			response := &pb.StopResponse{UnitId: id}

			if unit == nil {
				response.Error = fmt.Sprintf("unit %d not found", id)
				return stream.Send(response)
			}

			if unit.Status() == UnitStatusRunning {
				unit.Stop()
			}

			unit.Model.RestartsCount += 1
			db.Update(&unit.Model)

			unitCopy, err := s.StartUnit(unit.Model)

			if err != nil {
				response.Error = err.Error()
				response.Unit = unit.PB()
				return stream.Send(response)
			}

			response.Unit = unitCopy.PB()
			return stream.Send(response)
		})
	}

	return eg.Wait()
}

func (s *DaemonServer) deleteUnitsStream(
	unitIDs []uint64,
	stream pb.ProcessService_DeleteServer,
) error {
	db := s.Options.DBFactory()

	defer db.Commit()
	defer db.Close()

	eg, _ := errgroup.WithContext(stream.Context())

	for _, unitID := range unitIDs {
		id := unitID

		eg.Go(func() error {
			s.unitsMu.RLock()
			unit := s.units[id]
			s.unitsMu.RUnlock()

			if unit == nil {
				response := &pb.StopResponse{UnitId: id}
				response.Error = fmt.Sprintf("unit %d not found", id)
				return stream.Send(response)
			}

			unit.Stop()
			db.DeleteStruct(&unit.Model)

			s.unitsMu.Lock()
			delete(s.units, unit.Model.ID)
			s.unitsMu.Unlock()

			return stream.Send(&pb.StopResponse{
				UnitId: id,
				Unit:   unit.PB(),
			})
		})
	}

	return eg.Wait()
}

func (s *DaemonServer) filterUnitIDs(except []uint64) []uint64 {
	s.unitsMu.RLock()
	defer s.unitsMu.RUnlock()

	unitIDs := make([]uint64, 0, len(s.units))

unitLoop:
	for _, unit := range s.units {
		for _, exceptedUnitID := range except {
			if unit.Model.ID == exceptedUnitID {
				continue unitLoop
			}
		}

		unitIDs = append(unitIDs, unit.Model.ID)
	}

	return unitIDs
}

func (s *DaemonServer) Start(
	ctx context.Context,
	request *pb.StartRequest,
) (*pb.StartResponse, error) {
	db := s.Options.DBFactory()
	defer db.Close()

	tx, err := db.Begin(true)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer tx.Rollback()

	unitModel := UnitModel{
		Name: request.Name,
		Bin:  request.Bin,
		CWD:  request.Cwd,
		Args: request.Args,
	}

	if err := tx.Save(&unitModel); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	unit, err := s.StartUnit(unitModel)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := tx.Commit(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.StartResponse{
		Unit: unit.PB(),
	}

	return &response, nil
}

func (s *DaemonServer) List(context.Context, *emptypb.Empty) (*pb.ListResponse, error) {
	pbUnits := make([]*pb.Unit, len(s.units))
	idx := 0

	for _, unit := range s.units {
		pbUnits[idx] = unit.PB()
		idx++
	}

	response := pb.ListResponse{
		Units: pbUnits,
	}

	return &response, nil
}

func (s *DaemonServer) Stop(request *pb.StopRequest, stream pb.ProcessService_StopServer) error {
	return s.stopUnitsStream(request.UnitIds, stream)
}

func (s *DaemonServer) StopAll(
	request *pb.ExceptRequest,
	stream pb.ProcessService_StopAllServer,
) error {
	return s.stopUnitsStream(s.filterUnitIDs(request.UnitIds), stream)
}

func (s *DaemonServer) Restart(
	request *pb.StopRequest,
	stream pb.ProcessService_RestartServer,
) error {
	return s.restartUnitsStream(request.UnitIds, stream)
}

func (s *DaemonServer) RestartAll(
	request *pb.ExceptRequest,
	stream pb.ProcessService_RestartAllServer,
) error {
	return s.restartUnitsStream(s.filterUnitIDs(request.UnitIds), stream)
}

func (s *DaemonServer) Logs(request *pb.LogsRequest, stream pb.ProcessService_LogsServer) error {
	const (
		LogsTailDefault    uint64 = 32
		LogsTailMax        uint64 = 1000
		LogsMaxBytes       uint32 = 4 * 1024 * 1024 // 4 MB
		LogsFollowInterval        = time.Second
	)

	s.unitsMu.RLock()
	unit := s.units[request.UnitId]
	s.unitsMu.RUnlock()

	if unit == nil {
		return status.Errorf(codes.NotFound, "unit %d not found", request.UnitId)
	}

	logFilepath := s.getUnitLogFilepath(unit.Model.ID)
	logFile, err := os.OpenFile(logFilepath, os.O_RDONLY, logFilePerm)

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	defer logFile.Close()

	stat, err := logFile.Stat()

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	var tailAmount uint64

	if request.Lines == 0 {
		tailAmount = LogsTailDefault
	} else {
		tailAmount = min(request.Lines, LogsTailMax)
	}

	const NEWLINE byte = 0xA
	const RETURN byte = 0xD

	lines := make([]string, 0, tailAmount)

	var lineBytesCount uint32 = 0
	var lineBuf []byte = nil

	tailBreakPoint := -stat.Size()
	var tailCursor int64 = 0

	appendLine := func() {
		slices.Reverse(lineBuf)
		lines = append(lines, string(lineBuf))
	}

	for {
		tailCursor -= 1

		_, err := logFile.Seek(tailCursor, io.SeekEnd)

		if err != nil {
			break
		}

		charBuf := make([]byte, 1)
		_, err = logFile.Read(charBuf)

		if err != nil {
			break
		}

		char := charBuf[0]

		if char != NEWLINE && char != RETURN {
			lineBuf = append(lineBuf, char)
			lineBytesCount += 1

			if lineBytesCount >= LogsMaxBytes {
				appendLine()
				break
			}

			if tailCursor != tailBreakPoint {
				continue
			}
		}

		if len(lineBuf) > 0 {
			appendLine()
		}

		if uint64(len(lines)) >= tailAmount || tailCursor == tailBreakPoint {
			break
		}

		clear(lineBuf)
	}

	slices.Reverse(lines)

	err = stream.Send(&pb.LogsResponse{
		Lines: lines,
	})

	if err != nil || !request.Follow {
		return err
	}

	_, err = logFile.Seek(0, io.SeekEnd)

	if err != nil {
		return err
	}

	for {
		time.Sleep(LogsFollowInterval)
		scanner := bufio.NewScanner(logFile)

		var lines []string = nil

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if len(lines) == 0 {
			continue
		}

		if err := stream.Send(&pb.LogsResponse{Lines: lines}); err != nil {
			return err
		}
	}
}

func (s *DaemonServer) Delete(
	request *pb.StopRequest,
	stream pb.ProcessService_DeleteServer,
) error {
	return s.deleteUnitsStream(request.UnitIds, stream)
}

func (s *DaemonServer) DeleteAll(
	request *pb.ExceptRequest,
	stream pb.ProcessService_DeleteAllServer,
) error {
	return s.deleteUnitsStream(s.filterUnitIDs(request.UnitIds), stream)
}

func (s *DaemonServer) Show(
	ctx context.Context,
	request *pb.ShowRequest,
) (*pb.ShowResponse, error) {
	s.unitsMu.RLock()
	unit := s.units[request.UnitId]
	s.unitsMu.RUnlock()

	if unit == nil {
		return nil, status.Errorf(codes.NotFound, "unit %d not found", request.UnitId)
	}

	response := pb.ShowResponse{
		Id:      unit.Model.ID,
		Name:    unit.Model.Name,
		Cwd:     unit.Model.CWD,
		Command: strings.Join(unit.Command.Args, " "),
	}

	return &response, nil
}

func (s *DaemonServer) LogsClear(
	ctx context.Context,
	request *pb.LogsClearRequest,
) (*emptypb.Empty, error) {
	for _, id := range request.UnitIds {
		unitID := id

		if s.units[unitID] == nil {
			continue
		}

		go func() {
			logFilepath := s.getUnitLogFilepath(unitID)
			os.Truncate(logFilepath, 0)
		}()
	}

	return emptyResponse, nil
}

func createUnitStartCommand(ctx context.Context, model *UnitModel, logFile *os.File) *exec.Cmd {
	command := exec.CommandContext(ctx, model.Bin, model.Args...)
	command.Dir = model.CWD
	command.Stdout = logFile
	command.Stderr = logFile
	return command
}
