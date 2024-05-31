package daemon

import (
	"os"
	"os/exec"
)

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
