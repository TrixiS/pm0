package commands

import (
	"fmt"
	"os"

	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Start(ctx *command_context.CommandContext) error {
	if ctx.CLIContext.NArg() == 0 {
		return fmt.Errorf("specify a binary and optional args")
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	bin := ctx.CLIContext.Args().First()
	args := ctx.CLIContext.Args().Tail()

	request := pb.StartRequest{
		Cwd:  cwd,
		Bin:  bin,
		Args: args,
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Start(ctx.CLIContext.Context, &request)

		if err != nil {
			return err
		}

		fmt.Printf("started process with id %s and PID %d\n", response.Id, response.Pid)
		return nil
	})
}
