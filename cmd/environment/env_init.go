package environment

import (
	"context"
	stderrs "errors"
	"os"
	"path"
	"sync"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/internal/stdio/loader"
)

func newInitEnvCmd() *cobra.Command {
	constr := newEnvConstructor()
	c := &cobra.Command{
		Use:   "init",
		Short: "Init environment for projects in given folder",

		PreRunE: constr.preRun,
		RunE:    constr.runInit,

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	c.Flags().StringP(pathFlag, pathFlag[:1], "", `Path to folder with projects`)
	return c
}

func (c *envConstructor) runInit(cmd *cobra.Command, args []string) error {
	if c.checkIfEnvExist() {
		return c.askToRunTidy(cmd, args, "environment already exists")
	}

	progressChan := make(chan loader.Progress)
	gDone := loader.RunSeqLoader(context.Background(), c.io, progressChan)
	defer func() {
		<-gDone()
	}()

	defer func() {
		close(progressChan)
	}()

	var ldr loader.Progress
	{
		ldr = loader.NewInfiniteLoader("Initiating basis", loader.RectSpinner())
		progressChan <- ldr

		err := c.initBasis()
		if err != nil {
			ldr.Done(loader.DoneFailed)
			return errors.Wrap(err, "error initiating basis")
		}
		ldr.Done(loader.DoneSuccessful)
	}

	{
		ldr = loader.NewInfiniteLoader("Creating projects folders", loader.RectSpinner())
		progressChan <- ldr

		err := c.initProjectsDirs()
		if err != nil {
			ldr.Done(loader.DoneFailed)
			return errors.Wrap(err, "error initiating basis")
		}

		ldr.Done(loader.DoneSuccessful)
	}

	{
		ldr = loader.NewInfiniteLoader("Running rscli env tidy", loader.RectSpinner())
		progressChan <- ldr

		err := c.runTidy(cmd, args)
		if err != nil {
			ldr.Done(loader.DoneFailed)
			return errors.Wrap(err, "error initiating basis")
		}

		ldr.Done(loader.DoneSuccessful)
	}

	return nil
}

func (c *envConstructor) initBasis() error {
	err := stdio.CreateFolderIfNotExists(c.envDirPath)

	for _, f := range c.getSpirits() {
		err = stdio.CreateFileIfNotExists(path.Join(c.envDirPath, f.Name), f.Content)
		if err != nil {
			return errors.Wrap(err, "error creating file "+f.Name+" file")
		}
	}

	return nil
}
func (c *envConstructor) initProjectsDirs() error {
	wg := &sync.WaitGroup{}
	errC := make(chan error, len(c.srcProjDirs))

	wg.Add(len(c.srcProjDirs))
	for _, d := range c.srcProjDirs {
		go func(d os.DirEntry) {
			defer wg.Done()

			err := c.initProjectDir(d)
			if err != nil {
				errC <- errors.Wrap(err, "error creating "+d.Name())
			}

		}(d)
	}

	wg.Wait()

	close(errC)

	var errs []error
	for errP := range errC {
		errs = append(errs, errP)
	}

	if len(errs) == 0 {
		return nil
	}

	return stderrs.Join(errors.New("error preparing projects dirs"), stderrs.Join(errs...))
}
func (c *envConstructor) initProjectDir(d os.DirEntry) error {
	envProjDir := path.Join(c.envDirPath, d.Name())

	err := stdio.CreateFolderIfNotExists(envProjDir)
	if err != nil {
		return errors.Wrap(err, "error creating folder "+envProjDir)
	}

	var f []byte
	for _, spirit := range []patterns.File{c.selectMakefile(), patterns.EnvFile} {
		fPath := path.Join(c.envDirPath, spirit.Name)

		f, err = os.ReadFile(path.Join(c.envDirPath, spirit.Name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return errors.Wrap(err, "error reading "+fPath+" file")
		}

		err = stdio.CreateFileIfNotExists(path.Join(envProjDir, spirit.Name), f)
		if err != nil {
			return errors.Wrap(err, "error reading "+fPath+" file")
		}
	}

	return nil
}
