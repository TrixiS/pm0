package cli

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/TrixiS/pm0/internal/daemon/pb"
)

const allIdent string = "all"

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

	if allowAll && len(idents) == 1 && strings.ToLower(idents[0]) == allIdent {
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

func parseIDIdent(ident string) uint32 {
	intIdent, err := strconv.Atoi(ident)

	if err != nil {
		return 0 // 0 never equals to any unit model id
	}

	return uint32(intIdent)
}
