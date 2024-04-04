package daemon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"slices"
	"syscall"
	"time"

	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/asdine/storm/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const LogFileFlag = os.O_CREATE | os.O_RDWR | os.O_APPEND
const LogFilePerm = 0666

const LogsTailDefault uint64 = 32
const LogsTailMax uint64 = 1000

type DaemonServerOptions struct {
	LogsDirpath string
	DBFactory   func() *storm.DB
}

type DaemonServer struct {
	pb.UnimplementedProcessServiceServer

	Options DaemonServerOptions

	units map[UnitID]*Unit
}

func NewDaemonServer(options DaemonServerOptions) *DaemonServer {
	return &DaemonServer{
		Options: options,
		units:   make(map[UnitID]*Unit),
	}
}

func (s DaemonServer) getUnitLogFilepath(unitID UnitID) string {
	return path.Join(s.Options.LogsDirpath, fmt.Sprintf("%d.log", unitID))
}

func (s DaemonServer) openUnitLogFile(unitID UnitID) (*os.File, error) {
	logFilepath := s.getUnitLogFilepath(unitID)
	return os.OpenFile(logFilepath, LogFileFlag, LogFilePerm)
}

func (s DaemonServer) watchUnitProcess(unit *Unit) {
	err := unit.Command.Wait()

	unit.LogFile.Close()

	if err == nil {
		return
	}

	_, isExitErr := err.(*exec.ExitError)

	if !isExitErr {
		return
	}

	log.Printf("process %d exited: %v", unit.Model.ID, err)
}

func (s *DaemonServer) RestartUnit(model UnitModel) (*Unit, error) {
	logFile, err := s.openUnitLogFile(model.ID)

	if err != nil {
		return nil, err
	}

	command := createUnitStartCommand(model.Bin, model.Args, model.CWD, logFile)

	if err := command.Start(); err != nil {
		return nil, err
	}

	unitCopy := &Unit{
		Model:   model,
		Command: command,
		LogFile: logFile,
	}

	s.units[unitCopy.Model.ID] = unitCopy
	go s.watchUnitProcess(unitCopy)

	return unitCopy, nil
}

func (s *DaemonServer) Start(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {
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

	logFile, err := s.openUnitLogFile(unitModel.ID)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	command := createUnitStartCommand(request.Bin, request.Args, request.Cwd, logFile)

	if err = command.Start(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := tx.Commit(); err != nil {
		stopProcess(command.Process)
		command.Process.Release()
		return nil, status.Error(codes.Internal, err.Error())
	}

	unit := &Unit{
		Model:   unitModel,
		Command: command,
		LogFile: logFile,
	}

	s.units[unitModel.ID] = unit

	go s.watchUnitProcess(unit)

	response := pb.StartResponse{
		Unit: createUnitPB(unit),
	}

	return &response, nil
}

func (s *DaemonServer) List(context.Context, *emptypb.Empty) (*pb.ListResponse, error) {
	units := make([]*pb.Unit, len(s.units))
	unitIdx := 0

	for _, unit := range s.units {
		units[unitIdx] = createUnitPB(unit)
		unitIdx += 1
	}

	response := pb.ListResponse{
		Units: units,
	}

	return &response, nil
}

func (s *DaemonServer) Stop(request *pb.StopRequest, stream pb.ProcessService_StopServer) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, unitID := range request.UnitIds {
		id := UnitID(unitID)

		eg.Go(func() error {
			unit := s.units[id]
			response := &pb.StopResponse{UnitId: uint32(id)}

			if unit == nil {
				response.Error = fmt.Sprintf("unit %d not found", id)
				return stream.Send(response)
			}

			if unit.GetStatus() != RUNNING {
				response.Error = fmt.Sprintf("unit %s is not running", unit.Model.Name)
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			err := stopProcess(unit.Command.Process)

			if err != nil {
				response.Error = err.Error()
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			response.Unit = createUnitPB(unit)
			return stream.Send(response)
		})
	}

	return eg.Wait()
}

func (s *DaemonServer) Restart(request *pb.StopRequest, stream pb.ProcessService_RestartServer) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, unitID := range request.UnitIds {
		id := UnitID(unitID)

		eg.Go(func() error {
			unit := s.units[id]
			response := &pb.StopResponse{UnitId: uint32(id)}

			if unit == nil {
				response.Error = fmt.Sprintf("unit %d not found", id)
				return stream.Send(response)
			}

			unitStatus := unit.GetStatus()

			if unitStatus == RUNNING {
				if err := stopProcess(unit.Command.Process); err != nil {
					response.Error = err.Error()
					response.Unit = createUnitPB(unit)
					return stream.Send(response)
				}
			}

			unitCopy, err := s.RestartUnit(unit.Model)

			if err != nil {
				response.Error = err.Error()
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			s.units[unitCopy.Model.ID] = unitCopy

			go s.watchUnitProcess(unitCopy)

			response.Unit = createUnitPB(unitCopy)
			return stream.Send(response)
		})
	}

	return eg.Wait()
}

func (s *DaemonServer) Logs(request *pb.LogsRequest, stream pb.ProcessService_LogsServer) error {
	unit := s.units[UnitID(request.UnitId)]

	if unit == nil {
		return status.Error(codes.NotFound, "requested unit not found")
	}

	logFilepath := s.getUnitLogFilepath(unit.Model.ID)
	logFile, err := os.OpenFile(logFilepath, os.O_RDONLY, LogFilePerm)

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

	const NEWLINE byte = 10
	const RETURN byte = 13

	lines := make([]string, 0, tailAmount)
	var lineBuf []byte = nil

	tailBreakPoint := -stat.Size()
	var tailCursor int64 = 0

	for {
		tailCursor -= 1

		_, err := logFile.Seek(tailCursor, io.SeekEnd)

		if err != nil {
			break
		}

		hitBreakpoint := tailCursor == tailBreakPoint

		charBuf := make([]byte, 1)
		_, err = logFile.Read(charBuf)

		if err != nil {
			break
		}

		char := charBuf[0]

		if char != NEWLINE && char != RETURN {
			lineBuf = append(lineBuf, char)

			if !hitBreakpoint {
				continue
			}
		}

		if len(lineBuf) > 0 {
			slices.Reverse(lineBuf)
			lines = append(lines, string(lineBuf))
		}

		if uint64(len(lines)) >= tailAmount || hitBreakpoint {
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
		time.Sleep(time.Second)
		scanner := bufio.NewScanner(logFile)

		var lines []string = nil

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if len(lines) == 0 {
			continue
		}

		err = stream.Send(&pb.LogsResponse{
			Lines: lines,
		})

		if err != nil {
			return err
		}
	}
}

func stopProcess(process *os.Process) error {
	return process.Signal(syscall.SIGINT)
}

func createUnitPB(unit *Unit) *pb.Unit {
	var pid *int32
	var exitCode *int32

	unitStatus := unit.GetStatus()

	if unitStatus == RUNNING {
		int32Pid := int32(unit.Command.Process.Pid)
		pid = &int32Pid
	} else {
		int32ExitCode := int32(unit.Command.ProcessState.ExitCode())
		exitCode = &int32ExitCode
	}

	return &pb.Unit{
		Id:       uint32(unit.Model.ID),
		Name:     unit.Model.Name,
		Pid:      pid,
		Status:   uint32(unitStatus),
		ExitCode: exitCode,
	}
}

func createUnitStartCommand(
	bin string,
	args []string,
	cwd string,
	logFile *os.File,
) *exec.Cmd {
	command := exec.Command(bin, args...)
	command.Dir = cwd
	command.Stdout = logFile
	command.Stderr = logFile
	return command
}
