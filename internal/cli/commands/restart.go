package commands

import (
	"errors"
	"fmt"
	"io"

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

			err := stream.RecvMsg(&response)

			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}

			if response.Error != nil {
				fmt.Printf("failed to restart unit %d: %s\n", response.UnitId, *response.Error)
				continue
			}

			fmt.Printf("restarted unit %s with PID %d\n", response.Unit.Name, *response.Unit.Pid)
		}
	})
}
