package commands

import (
	"fmt"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func Rename(ctx *command.Context) error {
	args := ctx.CLI.Args()
	unitID, err := pm0.ParseStringUnitID(args.First())

	if err != nil {
		return err
	}

	name := args.Get(1)

	if name == "" {
		return fmt.Errorf("provide a new name with the second argument")
	}

	err = ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		_, err := client.Rename(ctx.CLI.Context, &pb.RenameRequest{
			UnitId: unitID,
			Name:   name,
		})

		return err
	})

	if err != nil {
		return err
	}

	pm0.Printf("renamed %s (%d)", name, unitID)
	return nil
}
