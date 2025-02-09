package commands

import (
	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Update(ctx *command.Context) error {
	unitID, err := pm0.ParseStringUnitID(ctx.CLI.Args().First())

	if err != nil {
		return err
	}

	name := ctx.CLI.String("name")

	err = ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		response, err := client.Update(ctx.CLI.Context, &pb.UpdateRequst{
			UnitId: unitID,
			Name:   name,
			Env:    ctx.CLI.StringSlice("env"),
		})

		if err != nil {
			return err
		}

		name = response.Name
		return nil
	})

	if err != nil {
		return err
	}

	pm0.Printf("updated %s (%d)", name, unitID)
	return nil
}
