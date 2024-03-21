package commands

import (
	"fmt"

	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var TableHeaderColorFunc = color.New(color.FgBlue, color.Underline).SprintfFunc()
var TableIDColorFunc = color.New(color.FgBlue, color.Bold).SprintfFunc()

func formatUnitStatus(status daemon.UnitStatus) string {
	switch status {
	case daemon.RUNNING:
		return "running"
	case daemon.EXITED:
		return "exited"
	case daemon.FAILED:
		return "failed"
	case daemon.STOPPED:
		return "stopped"
	default:
		return "unknown"
	}
}

func formatNillablePointer[T any](ptr *T) string {
	if ptr == nil {
		return "None"
	}

	return fmt.Sprintf("%v", *ptr)
}

func List(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.List(ctx.CLIContext.Context, nil)

		if err != nil {
			return err
		}

		headerColor := color.New(color.FgBlue, color.Underline).SprintfFunc()
		idColor := color.New(color.FgBlue, color.Bold).SprintfFunc()

		t := table.New("ID", "Name", "PID", "Status", "Exit code").
			WithHeaderFormatter(headerColor).
			WithFirstColumnFormatter(idColor)

		for _, unit := range response.Units {
			t.AddRow(
				unit.Id,
				unit.Name,
				formatNillablePointer(unit.Pid),
				formatUnitStatus(daemon.UnitStatus(unit.Status)),
				formatNillablePointer(unit.ExitCode),
			)
		}

		t.Print()
		return nil
	})
}
