package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
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
	ok, err := containsDependency(p.Cfg, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		// TODO: RSI-141
		// p.Io.Println("already contains pg dependency")
		return nil
	}

	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoClientFolderInConfig
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

	err = proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	return nil
}
