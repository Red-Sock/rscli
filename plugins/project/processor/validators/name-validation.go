package validators

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

var ErrInvalidNameErr = errors.New("name contains invalid symbol")

func ValidateName(p interfaces.Project) error {
	name := p.GetShortName()
	// starting and ending ascii symbols ranges that are applicable to project name
	availableRanges := [][]int32{
		{45, 47},
		{48, 57},
		{65, 90},
		{97, 122},
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
			return errors.Wrap(ErrInvalidNameErr, string(s))
		}
	}

	return nil
}
