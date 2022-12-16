package validators

import (
	"errors"
	"github.com/Red-Sock/rscli/plugins/src/project/processor/interfaces"
	"strings"
)

func ValidateName(p interfaces.Project) error {
	name := p.GetName()
	return ValidateNameString(name)
}

func ValidateNameString(name string) error {
	if name == "" {
		return errors.New("no name entered")
	}

	if strings.Contains(name, " ") {
		return errors.New("name contains space symbols")
	}

	return nil
}
