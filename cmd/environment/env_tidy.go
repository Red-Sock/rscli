package environment

import (
	"context"
	stderrs "errors"
	"path"

	"github.com/spf13/cobra"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/io/loader"
)

func newTidyEnvCmd() *cobra.Command {
	constr := newEnvConstructor()
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		RunE: constr.runTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `Path to folder with projects`)

	return c
}

func (c *envConstructor) runTidy(cmd *cobra.Command, arg []string) error {
	c.io.Println("Running rscli env tidy")

	err := c.initProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error during init of additional projects env dirs ")
	}

	err = c.fetchConstructor(cmd, arg)
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
		proj, err = project.LoadProjectEnvironment(c.cfg, path.Join(c.envDirPath, projName))
		if err != nil {
			return errors.Wrap(err, "error loading environment for project "+projName)
		}

		envPorts, err := proj.Environment.GetPortValues()
		if err != nil {
			return errors.Wrap(err, "error fetching ports for environment of "+projName)
		}

		for _, item := range envPorts {
			if conflictName := portManager.SaveIfNotExist(item.Value, item.Name); conflictName != "" {
				conflicts[item.Value] = []string{conflictName, item.Name}
			}
		}

		projectsEnvs[idx] = proj
	}

	done := loader.RunMultiLoader(context.Background(), c.io, progresses)
	defer func() {
		<-done()
		c.io.Println("rscli env tidy done")
	}()

	errC := make(chan error)
	for idx := range projectsEnvs {
		go func(i int) {
			err := projectsEnvs[i].Tidy(portManager, c.composePatterns)
			if err != nil {
				progresses[i].Done(loader.DoneFailed)
			} else {
				progresses[i].Done(loader.DoneSuccessful)
			}

			errC <- err
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
