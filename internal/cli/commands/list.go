package commands

import (
	"fmt"
	"time"

	"github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/rodaine/table"
)

func formatUnitStatus(status daemon.UnitStatus) string {
	switch status {
	case daemon.RUNNING:
		return "Running"
	case daemon.EXITED:
		return "Exited"
	case daemon.FAILED:
		return "Failed"
	case daemon.STOPPED:
		return "Stopped"
	default:
		return "Unknown"
	}
}

func formatNillablePointer[T any](ptr *T) string {
	if ptr == nil {
		return "None"
	}

	return fmt.Sprintf("%v", *ptr)
}

func getUnitUptime(startedAt int64, status daemon.UnitStatus) string {
	if status != daemon.RUNNING {
		return "None"
	}

	startedAtTime := time.Unix(startedAt, 0)
	now := time.Now()
	diff := now.Sub(startedAtTime)
	return diff.Round(time.Second).String()
}

func List(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.List(ctx.CLIContext.Context, nil)

		if err != nil {
			return err
		}

		t := table.New("ID", "Name", "PID", "Status", "Restarts", "Uptime").
			WithHeaderFormatter(cli.TableHeaderColorFunc).
			WithFirstColumnFormatter(cli.MainColorFunc)

		for _, unit := range response.Units {
			status := daemon.UnitStatus(unit.Status)

			t.AddRow(
				unit.Id,
				unit.Name,
				formatNillablePointer(unit.Pid),
				formatUnitStatus(status),
				unit.RestartsCount,
				getUnitUptime(unit.StartedAt, status),
			)
		}

		t.Print()
		return nil
	})
}
