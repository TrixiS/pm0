package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Stop(ctx *command_context.CommandContext) error {
	args := ctx.CLIContext.Args()

	if args.Len() == 0 {
		return errors.New("provide at least one unit identifier (id or name)")
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		request := &pb.StopRequest{
			Idents: args.Slice(),
		}

		stream, err := client.Stop(ctx.CLIContext.Context, request)

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

			if response.Success {
				fmt.Printf("stopped unit %s\n", response.Ident)
				continue
			}

			if response.Found {
				fmt.Printf("failed to stop unit %s\n", response.Ident)
				continue
			}

			fmt.Printf("running unit with ident %s not found\n", response.Ident)
		}
	})
}
