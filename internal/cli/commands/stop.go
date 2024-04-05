package commands

import (
	"context"
	"errors"
	"io"
	"strconv"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

func GetUnitIDsFromIdents(ctx context.Context, client pb.ProcessServiceClient, idents []string) ([]uint32, error) {
	response, err := client.List(ctx, nil)

	if err != nil {
		return nil, err
	}

	unitIDs := make([]uint32, 0, len(idents))

	for _, ident := range idents {
		intIdent, err := strconv.Atoi(ident)
		isIntIdent := true
		var unitID uint32 = 0

		if err == nil {
			unitID = uint32(intIdent)
		} else {
			isIntIdent = false
		}

		for _, unit := range response.Units {
			if unit.Name == ident || (isIntIdent && unit.Id == unitID) {
				unitIDs = append(unitIDs, unit.Id)
			}
		}
	}

	if len(unitIDs) == 0 {
		return unitIDs, errors.New("no units found for provided identifiers")
	}

	return unitIDs, nil
}

func Stop(ctx *command_context.CommandContext) error {
	args := ctx.CLIContext.Args()

	if args.Len() == 0 {
		return errors.New("provide at least one unit identifier (id or name)")
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := GetUnitIDsFromIdents(ctx.CLIContext.Context, client, args.Slice())

		if err != nil {
			return err
		}

		stream, err := client.Stop(ctx.CLIContext.Context, &pb.StopRequest{
			UnitIds: unitIDs,
		})

		if err != nil {
			return err
		}

		for {
			var response pb.StopResponse

			if err := stream.RecvMsg(&response); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}

			if response.Error == "" {
				pm0.Printf("stopped unit %s (%d)", response.Unit.Name, response.UnitId)
				continue
			}

			pm0.Printf("failed to stop unit %d: %s", response.UnitId, response.Error)
		}
	})
}
