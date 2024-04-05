package commands

import (
	"fmt"
	"os"

	pm0 "github.com/TrixiS/pm0/internal/cli"
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
		Name: ctx.CLIContext.String("name"),
		Bin:  bin,
		Args: args,
		Cwd:  cwd,
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Start(ctx.CLIContext.Context, &request)

		if err != nil {
			return err
		}

		pm0.Printf("started unit with id %d and PID %d", response.Unit.Id, *response.Unit.Pid)
		return nil
	})
}
