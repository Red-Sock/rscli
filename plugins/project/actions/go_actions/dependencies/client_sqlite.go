package dependencies

import (
	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/dockerfile_generator"
)

type Sqlite struct {
	dependencyBase
}

func sqlite(dep dependencyBase) Dependency {
	return &Sqlite{
		dep,
	}
}

func (s Sqlite) AppendToProject(proj Project) error {
	sc := sqlConn{Cfg: s.Cfg}

	err := sc.applySqlConnectionFile(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying changes to folder")
	}

	appNameInfo := proj.GetShortName()

	res := &resources.Sqlite{
		Name:             resources.SqliteResourceName,
		Path:             dockerfile_generator.DataVolumeName + "/" + appNameInfo + ".db",
		MigrationsFolder: "./migrations",
	}

	cfg := proj.GetConfig()
	if !containsDependency(cfg.DataSources, res) {
		cfg.DataSources = append(cfg.DataSources, res)
	}

	sc.applySqlDriver(proj, res.SqlDialect(), `_ "modernc.org/sqlite"`)

	return nil
}
