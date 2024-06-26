package commands

import (
	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func LogsClear(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := pm0.ParseUnitIDsFromArgs(ctx.CLIContext.Args().Slice())

		if err != nil {
			return err
		}

		_, err = client.LogsClear(ctx.CLIContext.Context, &pb.LogsClearRequest{UnitIds: unitIDs})

		if err != nil {
			return err
		}

		pm0.Printf("logs cleared")
		return nil
	})
}
