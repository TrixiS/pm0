package commands

import (
	"errors"
	"io"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Restart(ctx *command_context.CommandContext) error {
	args := ctx.CLIContext.Args()

	if args.Len() == 0 {
		return errors.New("provide at least one unit identifier (id or name)")
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := GetUnitIDsFromIdents(ctx.CLIContext.Context, client, args.Slice())

		if err != nil {
			return err
		}

		stream, err := client.Restart(ctx.CLIContext.Context, &pb.StopRequest{UnitIds: unitIDs})

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
				pm0.Printf(
					"restarted unit %s (%d) with PID %d",
					response.Unit.Name,
					response.Unit.Id,
					*response.Unit.Pid,
				)

				continue
			}

			pm0.Printf("failed to restart unit %d: %s", response.UnitId, response.Error)
		}
	})
}
