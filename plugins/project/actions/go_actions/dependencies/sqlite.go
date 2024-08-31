package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

type Sqlite struct {
	dependencyBase
}

func (s Sqlite) GetFolderName() string {
	if s.Name != "" {
		return s.Name
	}

	return resources.SqliteResourceName
}

func (s Sqlite) AppendToProject(proj project.Project) error {
	err := s.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying changes to folder")
	}

	s.applyConfig(proj)

	return nil
}

func (s Sqlite) applyClientFolder(proj project.Project) error {
	ok, err := containsDependencyFolder(s.Cfg.Env.PathsToClients, proj.GetFolder(), s.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	if len(s.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	sqliteConnFile := projpatterns.SqliteClientConnFile.
		CopyWithNewName(
			path.Join(
				s.Cfg.Env.PathsToClients[0],
				s.GetFolderName(),
				projpatterns.SqliteClientConnFile.Name,
			))

	proj.GetFolder().Add(sqliteConnFile)

	return nil
}

func (s Sqlite) applyConfig(proj project.Project) {
	for _, item := range proj.GetConfig().DataSources {
		if item.GetName() == s.GetFolderName() {
			return
		}
	}

	appNameInfo := proj.GetShortName()
	proj.GetConfig().DataSources = append(
		proj.GetConfig().DataSources,
		&resources.Sqlite{
			Name: resources.Name(s.GetFolderName()),
			Path: "./data/sqlite/" + appNameInfo + ".db",
		})
}
