package env

import (
	"context"
	stderrs "errors"
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/project"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/io/loader"
)

func (c *Constructor) RunTidy(cmd *cobra.Command, arg []string) error {
	c.Io.Println("Running rscli env tidy")

	err := c.initProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error during init of additional projects env dirs ")
	}

	err = c.FetchConstructor(cmd, arg)
	if err != nil {
		return errors.Wrap(err, "error fetching updated dirs")
	}

	portManager := ports.NewPortManager()

	progresses := make([]loader.Progress, len(c.envProjDirs))
	projectsEnvs := make([]*project.Env, len(c.envProjDirs))

	// TODO
	conflicts := make(map[uint16][]string)

	for idx := range c.envProjDirs {
		progresses[idx] = loader.NewInfiniteLoader(c.envProjDirs[idx].Name(), loader.RectSpinner())

		projName := c.envProjDirs[idx].Name()

		var proj *project.Env
		proj, err = project.LoadProjectEnvironment(c.Cfg, c.EnvManager.resources, path.Join(c.envDirPath, projName))
		if err != nil {
			return errors.Wrap(err, "error loading environment for project "+projName)
		}

		envPorts, err := proj.Environment.GetPortValues()
		if err != nil {
			return errors.Wrap(err, "error fetching ports for environment of "+projName)
		}

		for _, item := range envPorts {
			conflictServiceName := portManager.SaveIfNotExist(item.Value, item.Name)
			if conflictServiceName != "" {
				conflicts[item.Value] = []string{conflictServiceName, item.Name}
			}
		}
		proj.ComposePatterns = c.ComposePatterns

		projectsEnvs[idx] = proj
	}

	done := loader.RunMultiLoader(context.Background(), c.Io, progresses)
	defer func() {
		<-done()
		c.Io.Println("rscli env tidy done")
	}()

	errC := make(chan error)
	for idx := range projectsEnvs {
		go func(i int) {
			tidyErr := projectsEnvs[i].Tidy(portManager)
			if tidyErr != nil {
				progresses[i].Done(loader.DoneFailed)
			} else {
				progresses[i].Done(loader.DoneSuccessful)
			}

			errC <- tidyErr
		}(idx)
	}

	var errs []error
	for i := 0; i < len(c.envProjDirs); i++ {
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
