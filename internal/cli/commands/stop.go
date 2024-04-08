package commands

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
	"github.com/TrixiS/pm0/internal/daemon/pb"
)

const AllIdent string = "all"

func parseIDIdent(ident string) uint32 {
	intIdent, err := strconv.Atoi(ident)

	if err != nil {
		return 0 // 0 never equals to any unit model id
	}

	return uint32(intIdent)
}

func GetUnitIDsFromIdents(
	ctx context.Context,
	client pb.ProcessServiceClient,
	idents []string,
	allowAll bool,
) ([]uint32, error) {
	response, err := client.List(ctx, nil)

	if err != nil {
		return nil, err
	}

	if len(response.Units) == 0 {
		return nil, errors.New("unit list is empty")
	}

	unitIDs := make([]uint32, 0, len(idents))

	if allowAll && len(idents) == 1 && strings.ToLower(idents[0]) == AllIdent {
		for _, unit := range response.Units {
			unitIDs = append(unitIDs, unit.Id)
		}

		return unitIDs, nil
	}

	for _, ident := range idents {
		for _, unit := range response.Units {
			if unit.Name == ident || unit.Id == parseIDIdent(ident) {
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
		return ErrNoIdent
	}

	return ctx.Provider.WithClient(func(client pb.ProcessServiceClient) error {
		unitIDs, err := GetUnitIDsFromIdents(ctx.CLIContext.Context, client, args.Slice(), true)

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
