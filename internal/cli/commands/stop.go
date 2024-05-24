package commands

import (
	"errors"
	"io"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Stop(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := pm0.GetUnitIDsFromArgs(ctx.CLIContext, client, false)

		if err != nil {
			return err
		}

		stream, err := client.Stop(ctx.CLIContext.Context, &pb.StopRequest{
			UnitIds: unitIDs,
		})

		if err != nil {
			return err
		}

		for {
			var response pb.StopResponse

			if err := stream.RecvMsg(&response); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}

			if response.Error == "" {
				pm0.Printf("stopped unit %s (%d)", response.Unit.Name, response.UnitId)
				continue
			}

			pm0.Printf("failed to stop unit %d: %s", response.UnitId, response.Error)
		}
	})
}
