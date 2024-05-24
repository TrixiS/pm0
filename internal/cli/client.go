package cli

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/TrixiS/pm0/internal/daemon/pb"
	"github.com/urfave/cli/v2"
)

var ErrEmptyUnits = errors.New("unit list is empty")
var ErrNoUnitsFound = errors.New("no units found for provided identifiers")

const allIdent string = "all"

func GetUnitIDsFromArgs(ctx *cli.Context, client pb.ProcessServiceClient, single bool) ([]uint32, error) {
	args := ctx.Args()

	if args.Len() == 0 {
		return nil, ErrNoIdent
	}

	allowAll := !single

	var except []string = nil
	var idents []string = nil

	if allowAll {
		except = ctx.StringSlice("except")
		idents = args.Slice()
	} else {
		idents = []string{args.First()}
	}

	return getUnitIDsFromIdents(ctx.Context, client, idents, except, allowAll)
}

func getUnitIDsFromIdents(
	ctx context.Context,
	client pb.ProcessServiceClient,
	idents []string,
	except []string,
	allowAll bool,
) ([]uint32, error) {
	response, err := client.List(ctx, nil)

	if err != nil {
		return nil, err
	}

	if len(response.Units) == 0 {
		return nil, ErrEmptyUnits
	}

	unitIDs := make([]uint32, 0, len(idents))

	if allowAll && len(idents) > 0 && strings.ToLower(idents[0]) == allIdent {
	unitLoop:
		for _, unit := range response.Units {
			for _, exceptIdent := range except {
				if isUnitIdent(unit, exceptIdent) {
					continue unitLoop
				}
			}

			unitIDs = append(unitIDs, unit.Id)
		}

		return unitIDs, nil
	}

	for _, ident := range idents {
		for _, unit := range response.Units {
			if isUnitIdent(unit, ident) {
				unitIDs = append(unitIDs, unit.Id)
			}
		}
	}

	if len(unitIDs) == 0 {
		return unitIDs, ErrNoUnitsFound
	}

	return unitIDs, nil
}

func isUnitIdent(unit *pb.Unit, ident string) bool {
	return unit.Name == ident || unit.Id == parseIDIdent(ident)
}

func parseIDIdent(ident string) uint32 {
	intIdent, err := strconv.Atoi(ident)

	if err == nil {
		return uint32(intIdent)
	}

	return 0 // 0 never equals to any unit model id
}
