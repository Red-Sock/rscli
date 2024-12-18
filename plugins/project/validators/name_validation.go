package validators

import (
	"go.redsock.ru/rerrors"
)

var ErrInvalidNameErr = rerrors.New("name contains invalid symbol(s)")

func ValidateProjectNameStr(name string) error {
	// starting and ending ascii symbols ranges that are applicable to project name
	availableRanges := [][]int32{
		{45, 47},
		{48, 57},
		{65, 90},
		{97, 122},
		{95, 95},
	}
	for _, s := range name {
		var hasHitRange = false
		for _, r := range availableRanges {
			if s >= r[0] && s <= r[1] {
				hasHitRange = true
				break
			}
		}
		if !hasHitRange {
			return rerrors.Wrap(ErrInvalidNameErr, string(s))
		}
	}

	return nil
}
