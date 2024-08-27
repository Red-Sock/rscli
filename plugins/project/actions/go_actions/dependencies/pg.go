package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

type Postgres struct {
	Name string
	Cfg  *rscliconfig.RsCliConfig
	Io   io.StdIO
}

func (p Postgres) GetFolderName() string {
	if p.Name != "" {
		return p.Name
	}

	return resources.PostgresResourceName
}

func (p Postgres) AppendToProject(proj project.Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Postgres) applyClientFolder(proj project.Project) error {
	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	pgConnFile := projpatterns.PgConnFile.CopyWithNewName(
		path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), projpatterns.PgConnFile.Name))

	renamer.ReplaceProjectName(proj.GetName(), pgConnFile)

	proj.GetFolder().Add(
		pgConnFile,
		projpatterns.PgTxFile.CopyWithNewName(
			path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), projpatterns.PgTxFile.Name)),
	)

	return nil
}

func (p Postgres) applyConfig(proj project.Project) {
	for _, item := range proj.GetConfig().DataSources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}

	appNameInfo := proj.GetShortName()
	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, &resources.Postgres{
		Name:    resources.Name(p.GetFolderName()),
		Host:    "localhost",
		Port:    5432,
		DbName:  appNameInfo,
		User:    appNameInfo,
		Pwd:     "",
		SslMode: "",
	})
}
