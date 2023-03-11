package processor

import (
	"github.com/Red-Sock/rscli/plugins/project/processor/actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update"
	"github.com/pkg/errors"
)

func Tidy(pathToProject string) error {
	p, err := LoadProject(pathToProject)
	if err != nil {
		return err
	}

	err = actions.Tidy(p)
	if err != nil {
		return errors.Wrap(err, "error while tiding")
	}

	if p.RscliVersion.IsOlderThan(GetCurrentVersion()) {
		update.Do(p)
	}

	return nil
}
