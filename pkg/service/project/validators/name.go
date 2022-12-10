package validators

import (
	"errors"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/project/interfaces"
)

func ValidateName(p interfaces.Project) error {
	name := p.GetName()
	if name == "" {
		return errors.New("no name entered")
	}

	if strings.Contains(name, " ") {
		return errors.New("name contains space symbols")
	}
	return nil
}
