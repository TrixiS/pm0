package commands

import (
	"errors"
	"io"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Delete(ctx *command.Context) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := pm0.ParseUnitIDsFromArgs(ctx.CLI.Args().Slice())

		if err != nil {
			return err
		}

		stream, err := client.Delete(ctx.CLI.Context, &pb.StopRequest{UnitIds: unitIDs})

		if err != nil {
			return err
		}

		return readDeleteStream(stream)
	})
}

func DeleteAll(ctx *command.Context) error {
	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		stream, err := client.DeleteAll(
			ctx.CLI.Context,
			&pb.ExceptRequest{UnitIds: ctx.CLI.Uint64Slice("except")},
		)

		if err != nil {
			return err
		}

		return readDeleteStream(stream)
	})
}

func readDeleteStream(stream pb.ProcessService_DeleteClient) error {
	for {
		var response pb.StopResponse

		if err := stream.RecvMsg(&response); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		if response.Error == "" {
			pm0.Printf("deleted unit %s (%d)", response.Unit.Name, response.Unit.Id)
			continue
		}

		pm0.Printf("failed to delete unit %d: %s", response.UnitId, response.Error)
	}
}
