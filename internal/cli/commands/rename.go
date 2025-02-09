package commands

import (
	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Update(ctx *command.Context) error {
	name := ctx.CLI.String("name")
	unitID := ctx.CLI.Uint64("id")

	err := ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Update(ctx.CLI.Context, &pb.UpdateRequst{
			UnitId: unitID,
			Name:   name,
			Env:    ctx.CLI.StringSlice("env"),
		})

		name = response.Name
		return err
	})

	if err != nil {
		return err
	}

	pm0.Printf("updated %s (%d)", name, unitID)
	return nil
}
