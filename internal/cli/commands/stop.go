package commands

import (
	"errors"
	"io"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Stop(ctx *command.Context) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := pm0.ParseUnitIDsFromArgs(ctx.CLI.Args().Slice())

		if err != nil {
			return err
		}

		stream, err := client.Stop(ctx.CLI.Context, &pb.StopRequest{UnitIds: unitIDs})

		if err != nil {
			return err
		}

		return readStopStream(stream)
	})
}

func StopAll(ctx *command.Context) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		stream, err := client.StopAll(ctx.CLI.Context, &pb.ExceptRequest{
			UnitIds: ctx.CLI.Uint64Slice("except"),
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
