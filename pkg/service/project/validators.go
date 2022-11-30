package project

import (
	"errors"
	"strings"
)

func ValidateName(p *Project) error {
	if p.Name == "" {
		return errors.New("no name entered")
	}

	if strings.Contains(p.Name, " ") {
		return errors.New("name contains space symbols")
	}
	return nil
}
