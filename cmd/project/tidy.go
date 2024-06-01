package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
)

type projectTidy struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
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
		p.proj, err = project.LoadProject(p.path, p.config)
		if err != nil {
			return errors.Wrap(err, "error fetching project")
		}
	}

	err = go_actions.PrepareGoConfigFolderAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	err = go_actions.GenerateClientsAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error generating clients")
	}

	err = go_actions.GenerateServerAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error generating server")
	}

	err = go_actions.GenerateMakefileAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error generating makefiles")
	}

	err = go_actions.TidyAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error tiding project")
	}

	return nil
}
