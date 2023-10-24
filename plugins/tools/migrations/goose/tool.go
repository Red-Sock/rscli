package goose

import (
	"io"
	"net/http"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

const (
	installURL    = "github.com/pressly/goose/v3/cmd/goose"
	versionPrefix = "goose version: "
	versionURL    = "https://api.github.com/repos/pressly/goose/releases"
)

type Tool struct {
}

func (t *Tool) Install() error {
	_, err := cmd.Execute(cmd.Request{
		Tool: "go",
		Args: []string{"install", installURL},
	})
	if err != nil {
		return errors.Wrap(err, "error executing install command")
	}

	return nil
}

func (t *Tool) Version() (version string, err error) {
	res, err := cmd.Execute(cmd.Request{
		Tool: "goose",
		Args: []string{"-version"},
	})
	if err != nil {
		return "", errors.Wrap(err, "error executing goose version")
	}

	if !strings.HasPrefix(res, versionPrefix) {
		return res, errors.New("unexpected result of executing goose version")
	}

	return res, nil
}

func (t *Tool) GetLatestVersion() (version string, err error) {
	resp, err := http.Get(versionURL)
	if err != nil {
		return "", errors.Wrap(err, "error getting goose latest version")
	}

	var versionB []byte
	versionB, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response with goose latest version")
	}

	return string(versionB), nil
}
