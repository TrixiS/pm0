package daemon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/asdine/storm/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const LogFileFlag = os.O_CREATE | os.O_RDWR | os.O_APPEND
const LogFilePerm = 0666

type DaemonServerOptions struct {
	LogsDirpath string
	DBFactory   func() *storm.DB
}

type DaemonServer struct {
	pb.UnimplementedProcessServiceServer

	Options DaemonServerOptions

	units map[string]*Unit
}

func NewDaemonServer(options DaemonServerOptions) *DaemonServer {
	return &DaemonServer{
		Options: options,
		units:   make(map[string]*Unit),
	}
}

func (s DaemonServer) openUnitLogfile(unitID string) (*os.File, error) {
	logFilepath := path.Join(s.Options.LogsDirpath, fmt.Sprintf("%s.log", unitID))
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

	log.Printf("process %s exited: %v", unit.Model.ID, err)
}

func (s DaemonServer) getUnitByIdent(ident string) *Unit {
	unit := s.units[ident]

	if unit != nil {
		return unit
	}

	for _, unit := range s.units {
		if unit.Model.Name == ident {
			return unit
		}
	}

	return nil
}

func (s *DaemonServer) Start(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {
	unitID, err := gonanoid.New()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	logFile, err := s.openUnitLogfile(unitID)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	command := createUnitStartCommand(request.Bin, request.Args, request.Cwd, logFile)
	err = command.Start()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	unitModel := UnitModel{
		ID:   unitID,
		Name: request.Name,
		Bin:  request.Bin,
		CWD:  request.Cwd,
		Args: request.Args,
	}

	db := s.Options.DBFactory()
	err = db.Save(&unitModel)
	db.Close()

	if err != nil {
		stopProcess(command.Process)
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
		Id:  unitModel.ID,
		Pid: int32(command.Process.Pid),
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

	for _, ident := range request.Idents {
		i := ident

		eg.Go(func() error {
			unit := s.getUnitByIdent(i)
			response := &pb.StopResponse{Ident: i}

			if unit == nil {
				e := fmt.Sprintf("unit %s not found", i)
				response.Error = &e
				return stream.Send(response)
			}

			if unit.GetStatus() != RUNNING {
				e := fmt.Sprintf("unit %s is not running", unit.Model.Name)
				response.Error = &e
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			err := stopProcess(unit.Command.Process)

			if err != nil {
				e := err.Error()
				response.Error = &e
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

	for _, ident := range request.Idents {
		i := ident

		eg.Go(func() error {
			unit := s.getUnitByIdent(i)
			response := &pb.StopResponse{Ident: i}

			if unit == nil {
				e := fmt.Sprintf("unit %s not found", i)
				response.Error = &e
				return stream.Send(response)
			}

			unitStatus := unit.GetStatus()

			if unitStatus == RUNNING {
				if err := stopProcess(unit.Command.Process); err != nil {
					e := err.Error()
					response.Error = &e
					response.Unit = createUnitPB(unit)
					return stream.Send(response)
				}
			}

			logFile, err := s.openUnitLogfile(unit.Model.ID)

			if err != nil {
				e := err.Error()
				response.Error = &e
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			command := createUnitStartCommand(
				path.Base(unit.Command.Path),
				unit.Command.Args[1:],
				unit.Command.Dir,
				logFile,
			)

			err = command.Start()

			if err != nil {
				e := err.Error()
				response.Error = &e
				response.Unit = createUnitPB(unit)
				return stream.Send(response)
			}

			unitCopy := &Unit{
				Model:   unit.Model,
				Command: command,
				LogFile: logFile,
			}

			s.units[unitCopy.Model.ID] = unitCopy

			go s.watchUnitProcess(unitCopy)

			response.Unit = createUnitPB(unitCopy)
			return stream.Send(response)
		})
	}

	return eg.Wait()
}

func stopProcess(process *os.Process) error {
	if err := process.Signal(syscall.SIGINT); err != nil {
		return err
	}

	return process.Release()
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
		Id:       unit.Model.ID,
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
