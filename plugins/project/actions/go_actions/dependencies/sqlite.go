package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/plugins/project"
)

type Sqlite struct {
	dependencyBase
}

func (s Sqlite) GetFolderName() string {
	if s.Name != "" {
		return s.Name
	}

	return "sqldb"
}

func (s Sqlite) AppendToProject(proj project.Project) error {
	sc := sqlConn{Cfg: s.Cfg}

	err := sc.applySqlConnFile(proj)
	if err != nil {
		return errors.Wrap(err, "error applying changes to folder")
	}

	appNameInfo := proj.GetShortName()

	res := &resources.Sqlite{
		Name:             resources.SqliteResourceName,
		Path:             "./data/sqlite/" + appNameInfo + ".db",
		MigrationsFolder: "./migrations",
	}

	cfg := proj.GetConfig()
	if !containsDependency(cfg.DataSources, res) {
		cfg.DataSources = append(cfg.DataSources, res)
	}

	sc.applySqlDriver(proj, res.SqlDialect(), `_ "modernc.org/sqlite"`)

	return nil
}
