package git

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

func Init(workingDir string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    exe,
		Args:    []string{"init"},
		WorkDir: workingDir,
	})
	if err != nil {
		return errors.Wrap(err, "error initiating git repository")
	}

	return nil
}
