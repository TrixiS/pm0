package cli

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrEmptyUnits = errors.New("unit list is empty")

func ParseStringUnitID(id string) (uint64, error) {
	uint64ID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return uint64ID, fmt.Errorf("you should provide unit ids as unsigned integers: %w", err)
	}

	return uint64ID, nil
}

func ParseUnitIDsFromArgs(args []string) ([]uint64, error) {
	unitIDs := make([]uint64, len(args))

	for i, arg := range args {
		uint64ID, err := strconv.ParseUint(arg, 10, 64)

		if err != nil {
			return unitIDs, err
		}

		unitIDs[i] = uint64(uint64ID)
	}

	return unitIDs, nil
}
