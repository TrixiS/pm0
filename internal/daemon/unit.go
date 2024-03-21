package daemon

import (
	"os"
	"os/exec"
)

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
	LogFile *os.File
}

func (u Unit) GetStatus() UnitStatus {
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
