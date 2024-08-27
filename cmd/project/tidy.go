package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/go_project"
)

type projectTidy struct {
	io     io.IO
	config *config.RsCliConfig

	proj *go_project.Project
	path string
}

func newTidyCmd(pl projectTidy) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Cleans project",
		Long:  "Can be used clean project",

		RunE: pl.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// TODO
	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectTidy) run(_ *cobra.Command, _ []string) (err error) {
	if p.proj == nil {
		p.proj, err = go_project.LoadProject(p.path, p.config)
		if err != nil {
			return errors.Wrap(err, "error fetching project for tidy")
		}
	}

	ap := actions.NewActionPerformer(p.io, p.proj)

	err = ap.Tidy()
	if err != nil {
		return errors.Wrap(err, "error performing tidy")
	}

	return nil
}
