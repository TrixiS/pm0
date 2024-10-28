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
		unitIDs, err := pm0.ParseUnitIDsFromArgs(ctx.CLIContext.Args().Slice())

		if err != nil {
			return err
		}

		stream, err := client.Stop(ctx.CLIContext.Context, &pb.StopRequest{UnitIds: unitIDs})

		if err != nil {
			return err
		}

		return readStopStream(stream)
	})
}

func StopAll(ctx *command_context.CommandContext) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		stream, err := client.StopAll(ctx.CLIContext.Context, &pb.ExceptRequest{
			UnitIds: ctx.CLIContext.Uint64Slice("except"),
		})

		if err != nil {
			return err
		}

		return readStopStream(stream)
	})
}

func readStopStream(stream pb.ProcessService_StopClient) error {
	for {
		var response pb.StopResponse

		if err := stream.RecvMsg(&response); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		if len(response.Error) == 0 {
			pm0.Printf("stopped unit %s (%d)", response.Unit.Name, response.UnitId)
			continue
		}

		pm0.Printf("failed to stop unit %d: %s", response.UnitId, response.Error)
	}
}
