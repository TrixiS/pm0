package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
)

const setupScript string = `[Unit]
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

const daemonServiceFilepath string = "/etc/systemd/system/pm0_daemon.service"
const cliBinFilepath string = "/usr/bin/pm0"

func Setup(ctx *command_context.CommandContext) error {
	execFilepath, err := os.Executable()

	if err != nil {
		return err
	}

	execDirpath := path.Dir(execFilepath)
	cliFilepath := path.Join(execDirpath, "pm0")

	_, err = os.Stat(cliFilepath)

	if err != nil {
		return fmt.Errorf("pm0 cli executable not found")
	}

	serviceFile, err := os.OpenFile(daemonServiceFilepath, os.O_CREATE|os.O_RDWR, 0o770)

	if err != nil {
		return fmt.Errorf("open service file: %v", err)
	}

	defer serviceFile.Close()

	daemonFilepath := path.Join(execDirpath, "pm0_daemon")
	serviceString := fmt.Sprintf(setupScript, execDirpath, daemonFilepath)
	_, err = serviceFile.WriteString(serviceString)

	if err != nil {
		return fmt.Errorf("writing service file: %w", err)
	}

	err = exec.CommandContext(ctx.CLIContext.Context, "cp", cliFilepath, cliBinFilepath).Run()

	if err != nil {
		return fmt.Errorf("copy %s -> %s: %w", cliFilepath, cliBinFilepath, err)
	}

	err = exec.CommandContext(ctx.CLIContext.Context, "systemctl", "enable", daemonServiceFilepath).Run()

	if err != nil {
		return fmt.Errorf("failed to enable daemon service: %w", err)
	}

	err = exec.CommandContext(ctx.CLIContext.Context, "systemctl", "start", daemonServiceFilepath).Run()

	if err != nil {
		return fmt.Errorf("failed to start daemon service: %w", err)
	}

	pm0.Printf("pm0 daemon successfully set up")
	return nil
}
