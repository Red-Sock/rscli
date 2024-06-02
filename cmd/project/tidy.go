package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
)

type projectTidy struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
	path string
}

func tidySequence() []actions.Action {
	return []actions.Action{
		go_actions.PrepareGoConfigFolderAction{},
		go_actions.PrepareMakefileAction{},
		go_actions.PrepareClientsAction{},
		go_actions.BuildProjectAction{},
		go_actions.RunMakeGenAction{},
		go_actions.BuildProjectAction{},
		go_actions.RunGoTidyAction{},
	}
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

	err = tidy(p.io, p.proj)
	if err != nil {
		return errors.Wrap(err, "error performing tidy")
	}

	return nil
}

func tidy(printer io.IO, proj *project.Project) error {
	for _, a := range tidySequence() {
		printer.Println(a.NameInAction())
		err := a.Do(proj)
		if err != nil {
			return err
		}
	}

	return nil

}
