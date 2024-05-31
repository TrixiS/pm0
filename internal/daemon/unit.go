package daemon

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/TrixiS/pm0/internal/daemon/pb"
)

type UnitStatus uint32

const (
	RUNNING UnitStatus = 0
	EXITED  UnitStatus = 1
	FAILED  UnitStatus = 2
	STOPPED UnitStatus = 3
)

type UnitModel struct {
	ID            uint64 `storm:"id,increment"`
	Name          string
	CWD           string
	Bin           string
	Args          []string
	RestartsCount uint32
}

type Unit struct {
	Model     UnitModel
	Command   *exec.Cmd
	LogFile   *os.File
	StartedAt time.Time
	IsStopped bool
}

func (u Unit) GetStatus() UnitStatus {
	if u.IsStopped {
		return STOPPED
	}

	if u.Command.ProcessState == nil {
		return RUNNING
	}

	switch u.Command.ProcessState.ExitCode() {
	case -1:
		return STOPPED
	case 0:
		return EXITED
	default:
		return FAILED
	}
}

func (u Unit) ToPB() *pb.Unit {
	var pid *int32

	unitStatus := u.GetStatus()

	if unitStatus == RUNNING {
		int32Pid := int32(u.Command.Process.Pid)
		pid = &int32Pid
	}

	return &pb.Unit{
		Id:            u.Model.ID,
		Name:          u.Model.Name,
		Pid:           pid,
		Status:        uint32(unitStatus),
		RestartsCount: u.Model.RestartsCount,
		StartedAt:     u.StartedAt.Unix(),
	}
}

func (u *Unit) Stop() error {
	u.IsStopped = true
	return stopProcess(u.Command.Process)
}

func stopProcess(process *os.Process) error {
	return process.Signal(syscall.SIGINT)
}
