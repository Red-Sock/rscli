package environment

import (
	"context"
	stderrs "errors"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/loader"
	"github.com/Red-Sock/rscli/plugins/environment/env"
)

func newTidyEnvCmd(et *envTidy) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		PreRunE: et.constructor.FetchConstructor,
		RunE:    et.RunTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(env.PathFlag, env.PathFlag[:1], "", `Path to folder with projects`)
	c.Flags().BoolP(env.ServiceInContainer, env.ServiceInContainer[:1], false, "Service will be run in container")

	return c
}

type envTidy struct {
	io          io.IO
	constructor *env.Constructor
}

func (c *envTidy) RunTidy(cmd *cobra.Command, arg []string) error {
	c.io.Println("Running rscli env tidy")

	err := c.constructor.InitProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error during init of additional projects env dirs ")
	}

	tidyMngr, err := c.constructor.FetchTidyManager()
	if err != nil {
		return errors.Wrap(err, "error fetching tidy manager")
	}

	done := loader.RunMultiLoader(context.Background(), c.io, tidyMngr.Progresses)
	defer func() {
		<-done()
		c.io.Println("rscli env tidyMngr done")
	}()

	var serviceEnabled bool

	if cmd.Flag(env.ServiceInContainer).Value.String() == "true" {
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
	for i := 0; i < len(c.constructor.EnvProjDirs); i++ {
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
