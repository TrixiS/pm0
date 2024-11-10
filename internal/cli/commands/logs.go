package commands

import (
	"errors"
	"fmt"
	"io"
	"strings"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Logs(ctx *command.Context) error {
	unitID, err := pm0.ParseStringUnitID(ctx.CLI.Args().First())

	if err != nil {
		return err
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		stream, err := client.Logs(
			ctx.CLI.Context,
			&pb.LogsRequest{
				UnitId: unitID,
				Follow: ctx.CLI.Bool("follow"),
				Lines:  ctx.CLI.Uint64("lines"),
			},
		)

		if err != nil {
			return err
		}

		for {
			var response pb.LogsResponse

			err := stream.RecvMsg(&response)

			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}

			joinedLines := strings.Join(response.Lines, "\n")
			fmt.Println(joinedLines)
		}
	})
}
