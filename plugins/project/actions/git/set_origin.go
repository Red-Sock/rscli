package git

import (
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

func SetOrigin(wordDir, originURL string) error {
	if !strings.HasPrefix(originURL, "http") {
		originURL = "https://" + originURL
	}
	res, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"remote", "-v"},
		WorkDir: wordDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error listing remote repositories")
	}

	setRemote := cmd.Request{
		Tool:    bin,
		Args:    []string{"remote", "", "origin", originURL},
		WorkDir: wordDir,
	}

	if res == "" {
		setRemote.Args[1] = "add"
	} else {
		setRemote.Args[1] = "set-url"
	}

	_, err = cmd.Execute(setRemote)
	if err != nil {
		return rerrors.Wrap(err, "error setting origin url")
	}

	return nil
}
