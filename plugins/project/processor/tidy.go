package processor

import (
	"errors"

	errs "github.com/pkg/errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update"
)

var (
	ErrHasUncommittedChangesDuringTidy = errors.New("cannot execute tidy. Project has uncommitted changes")
)

func Tidy(pathToProject string) error {
	p, err := LoadProject(pathToProject)
	if err != nil {
		return err
	}

	status, err := actions.GitStatus(p)
	if err != nil {
		return errs.Wrap(err, "error while git status")
	}
	if len(status) != 0 {
		return errors.Join(ErrHasUncommittedChangesDuringTidy, errors.New(status.String()))
	}

	err = actions.Tidy(p)
	if err != nil {
		return errs.Wrap(err, "error while tiding")
	}

	if p.RscliVersion.IsOlderThan(GetCurrentVersion()) {
		err = update.Do(p)
		if err != nil {
			return err
		}
	}

	status, err = actions.GitStatus(p)
	if err != nil {
		return errs.Wrap(err, "error while getting git status after tidy")
	}
	if len(status) == 0 {
		return nil
	}

	err = actions.GitCommit(p.GetProjectPath(), "tidy commit:\n"+status.GetFilesListed())
	if err != nil {
		return errs.Wrap(err, "error executing commit")
	}

	return nil
}
