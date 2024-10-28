package daemon

import (
	"os"
	"os/exec"
	"time"

	"github.com/TrixiS/pm0/internal/daemon/pb"
)

type UnitStatus uint32

const (
	UnitStatusRunning UnitStatus = 0
	UnitStatusExited  UnitStatus = 1
	UnitStatusFailed  UnitStatus = 2
	UnitStatusStopped UnitStatus = 3
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
	Cancel    func()
}

func (u Unit) Status() UnitStatus {
	if u.Cancel == nil {
		return UnitStatusStopped
	}

	if u.Command.ProcessState == nil {
		return UnitStatusRunning
	}

	switch u.Command.ProcessState.ExitCode() {
	case -1:
		return UnitStatusStopped
	case 0:
		return UnitStatusExited
	default:
		return UnitStatusFailed
	}
}

func (u Unit) PB() *pb.Unit {
	var pid *int32

	unitStatus := u.Status()

	if unitStatus == UnitStatusRunning {
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

func (u *Unit) Stop() {
	if u.Cancel == nil {
		return
	}

	u.Cancel()
	u.Cancel = nil
}
