package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
)

const (
	serviceFileTemplate string = `[Unit]
Description="PM0 Daemon"

[Service]
User=root
Type=exec
KillSignal=SIGINT
WorkingDirectory=%s
ExecStart=%s
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target`

	daemonServiceFilename string = "pm0_daemon.service"
	daemonServiceFilepath string = "/etc/systemd/system/pm0_daemon.service"
	cliBinFilepath        string = "/usr/local/bin/pm0"
)

func Setup(ctx *command_context.CommandContext) error {
	exeFilepath, err := os.Executable()

	if err != nil {
		return fmt.Errorf("failed to get current executable filepath: %w", err)
	}

	if err := os.Symlink(exeFilepath, cliBinFilepath); err != nil {
		pm0.Printf(err.Error())
	} else {
		pm0.Printf("linked %s -> %s", exeFilepath, cliBinFilepath)
	}

	exeDirpath := path.Dir(exeFilepath)
	daemonFilepath := path.Join(exeDirpath, "pm0_daemon")

	if err := createServiceFile(exeDirpath, daemonFilepath); err != nil {
		return fmt.Errorf("create service file: %w", err)
	}

	err = exec.CommandContext(ctx.CLIContext.Context, "systemctl", "enable", daemonServiceFilename).
		Run()

	if err != nil {
		return fmt.Errorf("failed to enable daemon service: %w", err)
	}

	err = exec.CommandContext(ctx.CLIContext.Context, "systemctl", "start", daemonServiceFilename).
		Run()

	if err != nil {
		return fmt.Errorf("failed to start daemon service: %w", err)
	}

	pm0.Printf("pm0 daemon successfully set up")
	return nil
}

func createServiceFile(workingDirpath string, daemonExeFilepath string) error {
	serviceFile, err := os.OpenFile(daemonServiceFilepath, os.O_CREATE|os.O_RDWR, 0o770)

	if err != nil {
		return err
	}

	serviceString := fmt.Sprintf(serviceFileTemplate, workingDirpath, daemonExeFilepath)
	_, err = serviceFile.WriteString(serviceString)
	serviceFile.Close()
	return err
}
