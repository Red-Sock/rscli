package project

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/go_project"
)

const (
	pathFlag = "path"
)

type projectAdd struct {
	io     io.IO
	path   string
	config *rscliconfig.RsCliConfig

	proj *go_project.Project
}

func newAddCmd(projAdd projectAdd) *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds resource dependency to project",
		Long:  `Can be used to add a datasource or external API dependency to project`,

		RunE: projAdd.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectAdd) run(cmd *cobra.Command, args []string) error {
	err := p.loadProject(cmd)
	if err != nil {
		return errors.Wrap(err, "error loading project for add action")
	}

	p.io.Println("Searching for dependencies")

	for _, d := range dependencies.GetDependencies(p.config, args) {
		err = d.AppendToProject(p.proj)
		if err != nil {
			return errors.Wrap(err, "error adding dependency to project")
		}
	}

	p.io.Println("Dependencies added. Performing tidy")

	ap := actions.NewActionPerformer(p.io, p.proj)

	err = ap.Tidy()
	if err != nil {
		return errors.Wrap(err, "error tidying project")
	}

	p.io.Println("Tidy executed. Commiting changes")

	err = git.CommitWithUntracked(p.proj.GetProjectPath(), "added "+strings.Join(args, "; "))
	if err != nil {
		return errors.Wrap(err, "error performing git commit")
	}

	return nil
}

func (p *projectAdd) loadProject(cmd *cobra.Command) (err error) {
	var pathToProject string

	if cmd != nil {
		pathToProject = cmd.Flag(pathFlag).Value.String()
	}

	if pathToProject == "" {
		pathToProject = p.path
	}

	p.proj, err = go_project.LoadProject(pathToProject, p.config)
	if err != nil {
		return errors.Wrap(err, "error loading project")
	}

	return nil
}
