package commands

import (
	"os"

	pm0 "github.com/TrixiS/pm0/internal/cli"

	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func Show(ctx *command.Context) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitID, err := pm0.ParseStringUnitID(ctx.CLI.Args().First())

		if err != nil {
			return err
		}

		response, err := client.Show(ctx.CLI.Context, &pb.ShowRequest{UnitId: unitID})

		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)

		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Number: 1,
				Colors: text.Colors{text.Bold, text.FgHiCyan},
			},
			{
				Number: 2,
				Colors: text.Colors{text.FgHiWhite},
			},
		})

		t.AppendRows([]table.Row{
			{"ID", response.Id},
			{"Name", response.Name},
			{"CWD", response.Cwd},
			{"Command", response.Command},
		})

		t.Render()
		return nil
	})
}
