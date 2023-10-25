package goose

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
	"github.com/Red-Sock/rscli/plugins/tools/shared/ghversion"
)

const (
	installURL    = "github.com/pressly/goose/v3/cmd/goose@latest"
	versionPrefix = "goose version: "
	versionURL    = "https://api.github.com/repos/pressly/goose/releases/latest"
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

	var versionGH ghversion.GithubVersion
	err = json.Unmarshal(versionB, &versionGH)
	if err != nil {
		return "", errors.Wrap(err, "error unmarshalling gh response")
	}

	return versionGH.Name, nil
}

func (t *Tool) Migrate(pathToFolder string, resource resources.Resource) error {
	command := cmd.Request{
		Tool:    "goose",
		Args:    []string{string(resource.GetType()), "", "up"},
		WorkDir: pathToFolder,
	}

	envs := resource.GetEnv()

	switch resource.GetType() {
	case resources.DataSourcePostgres:
		command.Args[1] = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			envs[resources.EnvVarPostgresUser],
			envs[resources.EnvVarPostgresPassword],
			envs[resources.EnvVarPostgresHost],
			envs[resources.EnvVarPostgresPort],
			envs[resources.EnvVarPostgresDbName],
		)
	}

	res, err := cmd.Execute(command)
	if err != nil {
		return errors.Wrap(err, "error during migration")
	}

	// TODO
	_ = res
	return nil
}
