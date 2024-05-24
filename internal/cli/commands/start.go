package commands

import (
	"fmt"
	"os"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

// TODO: if name is empty string (make it optional) then decide name from cwd dirname
func Start(ctx *command_context.CommandContext) error {
	if ctx.CLIContext.NArg() == 0 {
		return fmt.Errorf("specify a binary and optional args")
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	name := ctx.CLIContext.String("name")

	if len(name) == 0 {
		name = path.Base(cwd)
	}

	bin := ctx.CLIContext.Args().First()
	args := ctx.CLIContext.Args().Tail()

	request := pb.StartRequest{
		Name: name,
		Bin:  bin,
		Args: args,
		Cwd:  cwd,
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Start(ctx.CLIContext.Context, &request)

		if err != nil {
			return err
		}

		pm0.Printf("started unit %s with id %d and PID %d", name, response.Unit.Id, *response.Unit.Pid)
		return nil
	})
}
