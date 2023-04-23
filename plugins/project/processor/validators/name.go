package validators

import (
	"errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

func ValidateName(p interfaces.Project) error {
	name := p.GetName()
	return ValidateNameString(name)
}

func ValidateNameString(name string) error {
	if name == "" {
		return errors.New("no name entered")
	}

	// starting and ending ascii symbols ranges that are applicable to project name
	availableRanges := [][]int32{
		{45, 45},
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
			return errors.New("name contains \"" + string(s) + "\" symbol")
		}
	}

	return nil
}
