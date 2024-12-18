package git

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

func Init(workingDir string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"init"},
		WorkDir: workingDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error initiating git repository")
	}

	return nil
}
