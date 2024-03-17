package daemon

import (
	"context"
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

	processLogsDirpath := path.Join(s.LogsDirpath, processID)

	if err := os.MkdirAll(processLogsDirpath, 0777); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stdoutFile, err := os.OpenFile(path.Join(processLogsDirpath, "out.log"), LogFileFlag, LogFilePerm)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stderrFile, err := os.OpenFile(path.Join(processLogsDirpath, "err.log"), LogFileFlag, LogFilePerm)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	command := exec.Command(request.Bin, request.Args...)
	command.Dir = request.Cwd
	command.Stdout = stdoutFile
	command.Stderr = stderrFile

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

		stdoutFile.Close()
		stderrFile.Close()

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
