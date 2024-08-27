package git

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/plugins/project"
)

type CommitWithUntrackedAction struct {
}

func (a CommitWithUntrackedAction) Do(p project.Project) error {
	err := CommitWithUntracked(p.GetProjectPath(), "rscli auto-commit")
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a CommitWithUntrackedAction) NameInAction() string {
	return "Commiting changes"
}

func Commit(workingDir, msg string) error {

	_, err := cmd.Execute(cmd.Request{
		Tool:    exe,
		Args:    []string{"commit", "-m", "\"" + msg + "\""},
		WorkDir: workingDir,
	})
	if err != nil {
		return errors.Wrap(err, "error committing files to git repository")
	}

	return nil

}

func CommitWithUntracked(workDir, msg string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    exe,
		Args:    []string{"add", "."},
		WorkDir: workDir,
	})
	if err != nil {
		return errors.Wrap(err, "error adding files to git repository")
	}

	status, err := Status(workDir)
	if err != nil {
		return errors.Wrap(err, "error getting git status")
	}
	if len(status) == 0 {
		return nil
	}

	err = Commit(workDir, msg)
	if err != nil {
		return errors.Wrap(err, "error performing commit")
	}

	return nil
}
