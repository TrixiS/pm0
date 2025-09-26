package commands

import (
	"fmt"
	"io"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Logs(ctx *command.Context) error {
	unitID, err := pm0.ParseStringUnitID(ctx.CLI.Args().First())

	if err != nil {
		return err
	}

	linesCount := ctx.CLI.Uint64("lines")
	follow := ctx.CLI.Bool("follow")

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		stream, err := client.Logs(
			ctx.CLI.Context,
			&pb.LogsRequest{
				UnitId: unitID,
				Follow: follow,
				Lines:  linesCount,
			},
		)

		if err != nil {
			return err
		}

		tailLines := make([]string, 0, linesCount)

		for {
			response := pb.LogsResponse{}

			if err := stream.RecvMsg(&response); err != nil {
				if err == io.EOF {
					break
				}

				return err
			}

			if !response.Flush {
				tailLines = append(tailLines, response.Line)
				continue
			}

			if len(tailLines) == 0 {
				fmt.Println(response.Line)
				continue
			}

			for i := len(tailLines) - 1; i >= 0; i-- {
				fmt.Println(tailLines[i])
			}

			tailLines = nil
		}

		return nil
	})
}
