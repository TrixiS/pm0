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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const LogFileFlag = os.O_CREATE | os.O_RDWR | os.O_APPEND
const LogFilePerm = 0666

type UnitStatus uint32

const (
	RUNNING UnitStatus = 0
	EXITED  UnitStatus = 1
	FAILED  UnitStatus = 2
	STOPPED UnitStatus = 3
)

type Unit struct {
	Model   UnitModel
	Command *exec.Cmd
}

func (u Unit) GetStatus() UnitStatus {
	if u.Command == nil {
		return STOPPED
	}

	if u.Command.ProcessState != nil {
		if u.Command.ProcessState.ExitCode() == 0 {
			return EXITED
		}

		return FAILED
	}

	return RUNNING
}

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
		command.Process.Signal(syscall.SIGINT)
		command.Process.Release()
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func() {
		err := command.Wait()

		logFile.Close()

		if err != nil {
			// TODO: hold process state and set it to exited on exit
			_, isExitErr := err.(*exec.ExitError)

			if !isExitErr {
				return
			}

			log.Printf("process %s exited: %v", processID, err)
		}
	}()

	s.units[processID] = &Unit{
		Model:   unitModel,
		Command: command,
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
			int32Pid := int32(unit.Command.Process.Pid)
			pid = &int32Pid
		} else {
			int32ExitCode := int32(unit.Command.ProcessState.ExitCode())
			exitCode = &int32ExitCode
		}

		units[unitIdx] = &pb.Unit{
			Id:       unitID,
			Name:     unit.Model.Name,
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
