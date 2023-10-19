package env

import (
	"context"
	stderrs "errors"
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/io/loader"
	"github.com/Red-Sock/rscli/plugins/environment/project"
	"github.com/Red-Sock/rscli/plugins/environment/project/ports"
)

type TidyManager struct {
	PortManager *ports.PortManager

	Progresses []loader.Progress
	ProjEnvs   []*project.ProjEnv

	conflicts map[uint16][]string
}

func (c *Constructor) RunTidy(cmd *cobra.Command, arg []string) error {
	c.io.Println("Running rscli env tidy")

	err := c.InitProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error during init of additional projects env dirs ")
	}

	err = c.FetchConstructor(cmd, arg)
	if err != nil {
		return errors.Wrap(err, "error fetching updated dirs")
	}

	tidyMngr, err := c.FetchTidyManager()
	if err != nil {
		return errors.Wrap(err, "error fetching tidy manager")
	}

	done := loader.RunMultiLoader(context.Background(), c.io, tidyMngr.Progresses)
	defer func() {
		<-done()
		c.io.Println("rscli env tidyMngr done")
	}()

	var serviceEnabled bool

	if cmd.Flag(ServiceInContainer).Value.String() == "true" {
		serviceEnabled = true
	}

	errC := make(chan error)
	for idx := range tidyMngr.ProjEnvs {
		go func(i int) {
			tidyErr := tidyMngr.ProjEnvs[i].Tidy(tidyMngr.PortManager, serviceEnabled)
			if tidyErr != nil {
				tidyMngr.Progresses[i].Done(loader.DoneFailed)
			} else {
				tidyMngr.Progresses[i].Done(loader.DoneSuccessful)
			}

			errC <- tidyErr
		}(idx)
	}

	var errs []error
	for i := 0; i < len(c.EnvProjDirs); i++ {
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

func (c *Constructor) FetchTidyManager() (*TidyManager, error) {
	var err error

	out := &TidyManager{
		PortManager: ports.NewPortManager(),

		Progresses: make([]loader.Progress, len(c.EnvProjDirs)),
		ProjEnvs:   make([]*project.ProjEnv, len(c.EnvProjDirs)),

		conflicts: make(map[uint16][]string),
	}

	for idx := range c.EnvProjDirs {
		out.Progresses[idx] = loader.NewInfiniteLoader(c.EnvProjDirs[idx].Name(), loader.RectSpinner())

		projName := c.EnvProjDirs[idx].Name()

		var proj *project.ProjEnv
		proj, err = project.LoadProjectEnvironment(
			c.cfg,
			c.environment,
			c.makefile,
			path.Join(c.envDirPath, projName),
			path.Join(path.Dir(c.envDirPath), projName),
		)
		if err != nil {
			return nil, errors.Wrap(err, "error loading environment for project "+projName)
		}

		envPorts, err := proj.Environment.GetPortValues()
		if err != nil {
			return nil, errors.Wrap(err, "error fetching ports for environment of "+projName)
		}

		for _, item := range envPorts {
			conflictServiceName := out.PortManager.SaveIfNotExist(item.Value, item.Name)
			if conflictServiceName != "" {
				out.conflicts[item.Value] = []string{conflictServiceName, item.Name}
			}
		}
		proj.ComposePatterns = c.composePatterns

		out.ProjEnvs[idx] = proj
	}

	return out, nil
}
