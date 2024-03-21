package daemon

import "os/exec"

type UnitStatus uint32

const (
	RUNNING UnitStatus = 0
	EXITED  UnitStatus = 1
	FAILED  UnitStatus = 2
	STOPPED UnitStatus = 3
)

type unit struct {
	model   UnitModel
	command *exec.Cmd
}

func (u unit) GetStatus() UnitStatus {
	if u.command.ProcessState == nil {
		return RUNNING
	}

	switch u.command.ProcessState.ExitCode() {
	case -1:
		return STOPPED
	case 0:
		return EXITED
	default:
		return FAILED
	}
}
