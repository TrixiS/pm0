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

	units map[string]*unit
}

func NewDaemonServer(options DaemonServerOptions) *DaemonServer {
	return &DaemonServer{
		Options: options,
		units:   make(map[string]*unit),
	}
}

// TODO: public func to start a unit (will be needed for resurrect)
func (s *DaemonServer) Start(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {
	processID, err := gonanoid.New()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	logFilepath := path.Join(s.Options.LogsDirpath, fmt.Sprintf("%s.log", processID))
	logFile, err := os.OpenFile(logFilepath, LogFileFlag, LogFilePerm)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	command := exec.Command(request.Bin, request.Args...)
	command.Dir = request.Cwd
	command.Stdout = logFile
	command.Stderr = logFile

	err = command.Start()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	unitModel := UnitModel{
		ID:   processID,
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

	go func() {
		err := command.Wait()

		logFile.Close()

		if err != nil {
			_, isExitErr := err.(*exec.ExitError)

			if !isExitErr {
				return
			}

			log.Printf("process %s exited: %v", processID, err)
		}
	}()

	s.units[processID] = &unit{
		model:   unitModel,
		command: command,
	}

	response := pb.StartResponse{
		Id:  unitModel.ID,
		Pid: int32(command.Process.Pid),
	}

	return &response, nil
}

func (s *DaemonServer) List(context.Context, *emptypb.Empty) (*pb.ListResponse, error) {
	units := make([]*pb.Unit, len(s.units))
	unitIdx := 0

	for unitID, unit := range s.units {
		var pid *int32
		var exitCode *int32

		unitStatus := unit.GetStatus()

		if unitStatus == RUNNING {
			int32Pid := int32(unit.command.Process.Pid)
			pid = &int32Pid
		} else {
			int32ExitCode := int32(unit.command.ProcessState.ExitCode())
			exitCode = &int32ExitCode
		}

		units[unitIdx] = &pb.Unit{
			Id:       unitID,
			Name:     unit.model.Name,
			Pid:      pid,
			Status:   uint32(unitStatus),
			ExitCode: exitCode,
		}

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
			isFound := unit != nil

			if !isFound || unit.command.ProcessState != nil {
				return stream.Send(&pb.StopResponse{
					Ident:   i,
					Found:   isFound,
					Success: false,
				})
			}

			err := stopProcess(unit.command.Process)

			if err != nil {
				return stream.Send(&pb.StopResponse{
					Ident:   i,
					Found:   true,
					Success: false,
				})
			}

			return stream.Send(&pb.StopResponse{
				Ident:   i,
				Found:   true,
				Success: true,
			})
		})
	}

	return eg.Wait()
}

func (s DaemonServer) getUnitByIdent(ident string) *unit {
	unit := s.units[ident]

	if unit != nil {
		return unit
	}

	for _, unit := range s.units {
		if unit.model.Name == ident && unit.GetStatus() == RUNNING {
			return unit
		}
	}

	return nil
}

func stopProcess(process *os.Process) error {
	if err := process.Signal(syscall.SIGINT); err != nil {
		return err
	}

	return process.Release()
}
