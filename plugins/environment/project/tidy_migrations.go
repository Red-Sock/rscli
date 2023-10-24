package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/project/config/resources"
	"github.com/Red-Sock/rscli/plugins/tools/migrations"
	"github.com/Red-Sock/rscli/plugins/tools/migrations/goose"
)

var (
	gooseMig  = &goose.Tool{}
	migrators = map[resources.DataSourceName]migrations.MigrationTool{
		resources.DataSourcePostgres: gooseMig,
	}
)

func (e *ProjEnv) tidyMigrations() error {
	migs, err := os.ReadDir(path.Join(e.pathToProjSrc, e.rscliConfig.Env.PathToMigrations))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return errors.Wrap(err, "error reading migrations dir")
	}

	ds, err := e.Config.GetDataSourceOptions()
	if err != nil {
		return errors.Wrap(err, "error getting datasource options")
	}

	for _, mig := range migs {
		if !mig.IsDir() {
			continue
		}

		var d resources.Resource
		for _, d = range ds {
			if d.GetName() == mig.Name() {
				break
			}
		}

		mt, ok := migrators[d.GetType()]
		if !ok {
			continue
		}
		_ = mt
		// TODO use migration tool
	}

	return nil
}
