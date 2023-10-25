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
	srcMigrationsFolder := path.Join(e.pathToProjSrc, e.rscliConfig.Env.PathToMigrations)
	targetMigrationFolder := path.Join(e.pathToProjInEnv, e.rscliConfig.Env.PathToMigrations)

	err := os.MkdirAll(targetMigrationFolder, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error creating migrations folder in env")
	}

	migs, err := os.ReadDir(srcMigrationsFolder)
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

	toolToMigrations := make(map[migrations.MigrationTool][]resources.Resource)

	for _, mig := range migs {
		if !mig.IsDir() {
			continue
		}

		var d resources.Resource
		for _, item := range ds {
			if item.GetName() == mig.Name() {
				d = item
				break
			}
		}

		if d == nil {
			continue
		}

		mt, ok := migrators[d.GetType()]
		if !ok {
			continue
		}

		toolToMigrations[mt] = append(toolToMigrations[mt], d)
	}

	var errs []error

	for tool, migPaths := range toolToMigrations {
		err = tool.Install()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var currentVersion, latestVersion string
		currentVersion, err = tool.Version()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		latestVersion, err = tool.GetLatestVersion()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if latestVersion > currentVersion {
			// TODO suggest to upgrade
		}

		for _, p := range migPaths {
			srcMigrations := path.Join(srcMigrationsFolder, p.GetName())
			targetMigrations := path.Join(targetMigrationFolder, p.GetName())
			err = os.RemoveAll(targetMigrations)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			err = os.Symlink(
				srcMigrations,
				targetMigrations,
			)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			err = tool.Migrate(srcMigrations, p)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	return nil
}
