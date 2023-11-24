package git

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

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

func ForceCommit(workDir, msg string) error {
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
