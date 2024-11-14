package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"
)

type Postgres struct {
	dependencyBase
}

func postgresClient(dep dependencyBase) Dependency {
	return &Postgres{
		dependencyBase: dep,
	}
}

func (p Postgres) AppendToProject(proj Project) error {
	sc := sqlConn{
		Cfg: p.dependencyBase.Cfg,
	}

	err := sc.applySqlConnFile(proj)
	if err != nil {
		return errors.Wrap(err, "error applying sql conn file")
	}

	appNameInfo := proj.GetShortName()

	res := &resources.Postgres{
		Name:             resources.PostgresResourceName,
		Host:             "localhost",
		Port:             5432,
		DbName:           appNameInfo,
		User:             appNameInfo,
		MigrationsFolder: "./migrations",
	}

	cfg := proj.GetConfig()
	if !containsDependency(cfg.DataSources, res) {
		cfg.DataSources = append(cfg.DataSources, res)
	}

	sc.applySqlDriver(proj, res.SqlDialect(), `_ "github.com/lib/pq"`)

	return nil
}