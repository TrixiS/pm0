package commands

import (
	"os"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func List(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.List(ctx.CLIContext.Context, nil)

		if err != nil {
			return err
		}

		if len(response.Units) == 0 {
			return pm0.ErrEmptyUnits
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
				pm0.FormatNillablePointer(unit.Pid),
				pm0.FormatUnitStatus(unitStatus),
				unit.RestartsCount,
				pm0.FormatUnitUptime(unit.StartedAt, unitStatus),
			})
		}

		t.Render()
		return nil
	})
}
