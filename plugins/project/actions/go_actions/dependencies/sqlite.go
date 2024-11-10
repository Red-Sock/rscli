package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"
)

type Sqlite struct {
	dependencyBase
}

func (s Sqlite) AppendToProject(proj Project) error {
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
