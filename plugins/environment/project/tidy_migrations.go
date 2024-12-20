package project

import (
	stderrs "errors"
	"os"
	"path"

	"go.redsock.ru/rerrors"
	"go.verv.tech/matreshka/resources"

	"github.com/Red-Sock/rscli/plugins/tools/migrations"
	"github.com/Red-Sock/rscli/plugins/tools/migrations/goose"
)

var (
	gooseMig  = &goose.Tool{}
	migrators = map[string]migrations.MigrationTool{
		resources.PostgresResourceName: gooseMig,
	}
)

func (e *ProjEnv) tidyMigrationDirs() error {
	srcMigrationsFolder := path.Join(e.pathToProjSrc, e.rscliConfig.Env.PathToMigrations)
	targetMigrationFolder := path.Join(e.pathToProjInEnv, e.rscliConfig.Env.PathToMigrations)

	migs, err := os.ReadDir(srcMigrationsFolder)
	if err != nil {
		if rerrors.Is(err, os.ErrNotExist) {
			return nil
		}

		return rerrors.Wrap(err, "error reading migrations dir")
	}

	err = os.MkdirAll(targetMigrationFolder, os.ModePerm)
	if err != nil {
		return rerrors.Wrap(err, "error creating migrations folder in env")
	}

	toolToMigrations := make(map[migrations.MigrationTool][]resources.Resource)

	for _, mig := range migs {
		if !mig.IsDir() {
			continue
		}

		var d resources.Resource
		for _, item := range e.Config.DataSources {
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
		}
	}

	return stderrs.Join(errs...)
}
