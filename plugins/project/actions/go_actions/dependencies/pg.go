package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type Postgres struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (Postgres) GetFolderName() string {
	return resources.PostgresResourceName
}

func (p Postgres) AppendToProject(proj interfaces.Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Postgres) applyClientFolder(proj interfaces.Project) error {
	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		return nil
	}

	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	pgConnFile := projpatterns.PgConnFile.CopyWithNewName(
		path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), projpatterns.PgConnFile.Name))

	go_actions.ReplaceProjectName(proj.GetName(), pgConnFile)

	proj.GetFolder().Add(
		pgConnFile,
		projpatterns.PgTxFile.CopyWithNewName(
			path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), projpatterns.PgTxFile.Name)),
	)

	return nil
}

func (p Postgres) applyConfig(proj interfaces.Project) {
	for _, item := range proj.GetConfig().Resources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}
	appNameInfo := proj.GetConfig().AppInfo.Name
	proj.GetConfig().Resources = append(proj.GetConfig().Resources, &resources.Postgres{
		Name:    resources.Name(p.GetFolderName()),
		Host:    "localhost",
		Port:    5432,
		DbName:  appNameInfo,
		User:    appNameInfo,
		Pwd:     "",
		SSLMode: "",
	})
}
