package commands

import (
	"fmt"
	"os"
	"time"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var runningStatusString = text.FgGreen.Sprint("Running")
var exitedStatusString = text.FgWhite.Sprint("Exited")
var failedStatusString = text.FgRed.Sprint("Failed")
var stoppedStatusString = text.FgYellow.Sprint("Stopped")

const tableNoneString = "None"

func formatNillablePointer[T any](ptr *T) string {
	if ptr == nil {
		return tableNoneString
	}

	return fmt.Sprintf("%v", *ptr)
}

func getUnitUptime(startedAt int64, status daemon.UnitStatus) string {
	if status != daemon.RUNNING {
		return tableNoneString
	}

	startedAtTime := time.Unix(startedAt, 0)
	now := time.Now()
	diff := now.Sub(startedAtTime)
	return diff.Round(time.Second).String()
}

func formatUnitStatus(unitStatus daemon.UnitStatus) string {
	switch unitStatus {
	case daemon.RUNNING:
		return runningStatusString
	case daemon.EXITED:
		return exitedStatusString
	case daemon.FAILED:
		return failedStatusString
	case daemon.STOPPED:
		return stoppedStatusString
	default:
		return "Unknown"
	}
}

func List(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.List(ctx.CLIContext.Context, nil)

		if err != nil {
			return err
		}

		if len(response.Units) == 0 {
			pm0.Printf("no units")
			return nil
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Name", "PID", "Status", "Restarts", "Uptime"})
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = false
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:   "ID",
				Colors: text.Colors{text.Bold, text.FgHiCyan},
			},
		})

		for _, unit := range response.Units {
			unitStatus := daemon.UnitStatus(unit.Status)

			t.AppendRow(table.Row{
				unit.Id,
				unit.Name,
				formatNillablePointer(unit.Pid),
				formatUnitStatus(unitStatus),
				unit.RestartsCount,
				getUnitUptime(unit.StartedAt, unitStatus),
			})
		}

		t.Render()
		return nil
	})
}
