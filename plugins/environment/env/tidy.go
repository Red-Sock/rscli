package env

import (
	"context"
	stderrs "errors"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/loader"
	"github.com/Red-Sock/rscli/plugins/environment/project"
	"github.com/Red-Sock/rscli/plugins/environment/project/ports"
)

func (e *GlobalEnvironment) Tidy() error {

	err := e.fetchFiles()
	if err != nil {
		return errors.Wrap(err, "error fetching global environment")
	}

	progresses, projEnvs, err := e.collectProjectEnvironments()
	if err != nil {
		return errors.Wrap(err, "error collecting environment info")
	}

	return e.run(progresses, projEnvs)
}

func (e *GlobalEnvironment) collectProjectEnvironments() ([]loader.Progress, []*project.ProjEnv, error) {
	portManager := ports.NewPortManager()

	progresses := make([]loader.Progress, len(e.envProjDirs))
	projEnvs := make([]*project.ProjEnv, len(e.envProjDirs))
	conflicts := make(map[uint16][]string)

	for idx := range e.envProjDirs {
		progresses[idx] = loader.NewInfiniteLoader(e.envProjDirs[idx].Name(), loader.RectSpinner())

		projName := e.envProjDirs[idx].Name()

		proj, err := project.LoadProjectEnvironment(
			e.rsCliConfig,
			e.environment,
			e.makefile,
			e.composePatterns,

			portManager,

			path.Join(e.envDirPath, projName),
			path.Join(path.Dir(e.envDirPath), projName),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error loading environment for project "+projName)
		}

		envPorts, err := proj.Environment.GetPortValues()
		if err != nil {
			return nil, nil, errors.Wrap(err, "error fetching ports for environment of "+projName)
		}

		for _, item := range envPorts {
			conflictServiceName := portManager.SaveIfNotExist(item.Value, item.Name)
			if conflictServiceName != "" {
				conflicts[item.Value] = []string{conflictServiceName, item.Name}
			}
		}

		projEnvs[idx] = proj
	}
	return progresses, projEnvs, nil
}

func (e *GlobalEnvironment) run(progresses []loader.Progress, envs []*project.ProjEnv) error {
	done := loader.RunMultiLoader(context.Background(), e.io, progresses)
	defer func() {
		<-done()
		e.io.Println("rscli env tidyMngr done")
	}()

	errC := make(chan error)
	for idx := range envs {
		go func(i int) {
			// TODO
			tidyErr := envs[i].Tidy(false)
			if tidyErr != nil {
				progresses[i].Done(loader.DoneFailed)
			} else {
				progresses[i].Done(loader.DoneSuccessful)
			}

			errC <- tidyErr
		}(idx)
	}

	var errs []error
	for i := 0; i < len(e.envProjDirs); i++ {
		err, ok := <-errC
		if !ok {
			break
		}

		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}

	return stderrs.Join(errs...)
}
