package environment

import (
	"context"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/env"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/internal/io/loader"
)

func newInitEnvCmd(ei *envInit) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Init environment for projects in given folder",

		PreRunE: ei.constructor.FetchConstructor,
		RunE:    ei.RunInit,

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	c.Flags().StringP(env.PathFlag, env.PathFlag[:1], "", `Path to folder with projects`)
	return c
}

type envInit struct {
	io          io.IO
	constructor *env.Constructor
}

func (c *envInit) RunInit(cmd *cobra.Command, args []string) error {
	if c.constructor.IsEnvExist() {
		return c.askToRunTidy(cmd, args, "environment already exists", colors.ColorYellow)
	}

	err := c.runInit()
	if err != nil {
		return errors.Wrap(err, "error during basic init ")
	}

	return c.askToRunTidy(cmd, args, "environment prepared", colors.ColorGreen)
}

func (c *envInit) runInit() error {
	progressChan := make(chan loader.Progress)
	gDone := loader.RunSeqLoader(context.Background(), c.io, progressChan)

	defer func() { <-gDone() }()
	defer func() { close(progressChan) }()

	var ldr loader.Progress
	{
		ldr = loader.NewInfiniteLoader("Initiating basis", loader.RectSpinner())
		progressChan <- ldr

		err := c.constructor.InitBasis()
		if err != nil {
			ldr.Done(loader.DoneFailed)
			return errors.Wrap(err, "error initiating basis")
		}
		ldr.Done(loader.DoneSuccessful)
	}

	{
		ldr = loader.NewInfiniteLoader("Creating projects folders", loader.RectSpinner())
		progressChan <- ldr

		err := c.constructor.InitProjectsDirs()
		if err != nil {
			ldr.Done(loader.DoneFailed)
			return errors.Wrap(err, "error initiating basis")
		}

		ldr.Done(loader.DoneSuccessful)
	}

	return nil
}

func (c *envInit) askToRunTidy(cmd *cobra.Command, args []string, msg string, color colors.Color) error {
	c.io.Println()
	c.io.PrintColored(color, msg+
		"!\nWant to run \"rscli env tidyManager\"? (Y)es/(N)o: ")

	for {
		resp, err := c.io.GetInput()
		if err != nil {
			return errors.Wrap(err, "error obtaining user input")
		}
		r := strings.ToLower(resp)[0]
		if r == 'y' {
			et := envTidy{io: c.io, constructor: c.constructor}
			return et.RunTidy(cmd, args)
		}

		if r == 'n' {
			return nil
		}
		c.io.PrintlnColored(colors.ColorRed, "No can't do "+resp+"! \"(Y)es\" or \"(N)o\":")
	}

	return nil
}
