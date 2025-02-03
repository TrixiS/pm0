package commands

import (
	"fmt"
	"os"
	"path"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Start(ctx *command.Context) error {
	if ctx.CLI.NArg() == 0 {
		return fmt.Errorf("specify a binary and optional args")
	}

	osCwd, err := os.Getwd()

	if err != nil {
		return err
	}

	var cwd string

	if ctx.CLI.IsSet("cwd") {
		cwd = ctx.CLI.String("cwd")
	} else {
		cwd = osCwd
	}

	name := ctx.CLI.String("name")

	if len(name) == 0 {
		name = path.Base(osCwd)
	}

	bin := ctx.CLI.Args().First()
	args := ctx.CLI.Args().Tail()

	request := pb.StartRequest{
		Name: name,
		Bin:  bin,
		Args: args,
		Cwd:  cwd,
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Start(ctx.CLI.Context, &request)

		if err != nil {
			return err
		}

		pm0.Printf("started unit %s (%d) with PID %d", name, response.Unit.Id, *response.Unit.Pid)
		return nil
	})
}
