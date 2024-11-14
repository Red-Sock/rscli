package git

import (
	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/project"
)

const (
	ChangesTypeInvalid = iota
	ChangesTypeNotStaged
	ChangesTypeNotCommitted
)

const bin = "git"

type InitGit struct{}

func (a InitGit) Do(p project.IProject) error {
	projectPath := p.GetProjectPath()

	err := Init(projectPath)
	if err != nil {
		return errors.Wrap(err, "error initializing project")
	}

	err = CommitWithUntracked(projectPath, "project init via RedSock CLI")
	if err != nil {
		return errors.Wrap(err, "error committing changes")
	}

	err = SetOrigin(projectPath, p.GetName())
	if err != nil {
		return errors.Wrap(err, "error setting git origin")
	}

	return nil
}
func (a InitGit) NameInAction() string {
	return "Initiating git"
}
