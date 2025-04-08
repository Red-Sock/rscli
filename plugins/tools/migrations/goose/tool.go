package goose

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/plugins/tools/shared/ghversion"
)

const (
	installURL    = "github.com/pressly/goose/v3/cmd/goose@latest"
	versionPrefix = "goose version: "
	versionURL    = "https://api.github.com/repos/pressly/goose/releases/latest"
)

var ErrUnknownResourceToMigrate = rerrors.New("unknown resource to perform migration")

type Tool struct {
}

func (t *Tool) Install() error {
	_, err := cmd.Execute(cmd.Request{
		Tool: "go",
		Args: []string{"install", installURL},
	})
	if err != nil {
		return rerrors.Wrap(err, "error executing install command")
	}

	return nil
}

func (t *Tool) Version() (version string, err error) {
	res, err := cmd.Execute(cmd.Request{
		Tool: "goose",
		Args: []string{"-version"},
	})
	if err != nil {
		return "", rerrors.Wrap(err, "error executing goose version")
	}

	if !strings.HasPrefix(res, versionPrefix) {
		return res, rerrors.New("unexpected result of executing goose version")
	}

	return res, nil
}

func (t *Tool) GetLatestVersion() (version string, err error) {
	resp, err := http.Get(versionURL)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting goose latest version")
	}

	var versionB []byte
	versionB, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading response with goose latest version")
	}

	var versionGH ghversion.GithubVersion
	err = json.Unmarshal(versionB, &versionGH)
	if err != nil {
		return "", rerrors.Wrap(err, "error unmarshalling gh response")
	}

	return versionGH.Name, nil
}

func (t *Tool) Migrate(pathToFolder string, resource resources.Resource) error {

	switch resource.GetType() {
	case resources.PostgresResourceName:
		return t.MigratePostgres(pathToFolder, resource)
	default:
		return ErrUnknownResourceToMigrate
	}

}

func (t *Tool) MigratePostgres(pathToFolder string, resource resources.Resource) error {
	command := cmd.Request{
		Tool:    "goose",
		Args:    []string{resource.GetType(), "", "up"},
		WorkDir: pathToFolder,
	}

	pg, ok := resource.(*resources.Postgres)
	if !ok {
		return rerrors.Wrapf(matreshka.ErrUnexpectedType, "expected postgres, got %T", resource)
	}

	command.Args[1] = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		pg.User,
		pg.Pwd,
		pg.Host,
		pg.Port,
		pg.DbName,
	)

	res, err := cmd.Execute(command)
	if err != nil {
		return rerrors.Wrap(err, "error during migration")
	}

	// Todo
	_ = res

	return nil
}
