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
)

const LogFileFlag = os.O_CREATE | os.O_RDWR | os.O_APPEND
const LogFilePerm = 0666

type DaemonServer struct {
	pb.UnimplementedProcessServiceServer

	LogsDirpath string
	DBFactory   func() *storm.DB
}

func (s *DaemonServer) Start(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {
	processID, err := gonanoid.New()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	logFilepath := path.Join(s.LogsDirpath, fmt.Sprintf("%s.log", processID))
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

	process := Process{
		ID:   processID,
		Bin:  request.Bin,
		CWD:  request.Cwd,
		Args: request.Args,
		PID:  command.Process.Pid,
	}

	db := s.DBFactory()
	err = db.Save(&process)
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

	response := pb.StartResponse{
		Id:  process.ID,
		Pid: int32(process.PID),
	}

	return &response, nil
}
