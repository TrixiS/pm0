package commands

import (
	"os"

	pm0 "github.com/TrixiS/pm0/internal/cli"

	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func Show(ctx *command_context.CommandContext) error {
	args := ctx.CLIContext.Args()

	if args.Len() == 0 {
		return pm0.ErrNoIdent
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := pm0.GetUnitIDsFromIdents(ctx.CLIContext.Context, client, []string{args.First()}, false)

		if err != nil {
			return err
		}

		response, err := client.Show(ctx.CLIContext.Context, &pb.ShowRequest{UnitId: unitIDs[0]})

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
