package project

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
)

type projectRun struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
	path string
}

func newRun(pr projectRun) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Run project and it's dependencies",
		Long:  "Can be used to prepare dev environment",

		RunE: pr.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (pr *projectRun) run(cmd *cobra.Command, args []string) error {
	// TODO
	return nil
}
