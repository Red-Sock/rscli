package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type Postgres struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (Postgres) GetFolderName() string {
	return "postgres"
}

func (p Postgres) Do(proj interfaces.Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Postgres) applyClientFolder(proj interfaces.Project) error {
	ok, err := containsDependency(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		return nil
	}

	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.ConnFileName),
			Content: patterns.PgConnFile,
		},
		&folder.Folder{
			Name:    path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.PgTxFileName),
			Content: patterns.PgTxFile,
		},
	)

	return nil
}

func (p Postgres) applyConfig(proj interfaces.Project) {
	ds := proj.GetConfig().DataSources
	if _, ok := ds[p.GetFolderName()]; ok {
		return
	}

	ds[p.GetFolderName()] = resources.Postgres{
		ResourceName: p.GetFolderName(),
		Host:         "localhost",
		Port:         5432,
		Name:         "",
		User:         "",
		Pwd:          "",
		SSLMode:      "",
	}
}
