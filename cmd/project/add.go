package project

import (
	"os"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
)

const (
	pathFlag = "path"
)

type projectAdd struct {
	io     io.IO
	path   string
	config *rscliconfig.RsCliConfig

	proj *project.Project
}

func newAddCmd(projAdd projectAdd) *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds dependency to project project",
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
		return errors.Wrap(err, "error loading project")
	}

	deps := dependencies.GetDependencies(p.config, args)
	for _, d := range deps {
		err = d.AppendToProject(p.proj)
		if err != nil {
			return errors.Wrap(err, "error adding dependency to project")
		}
	}

	err = p.proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building folders")
	}

	b, err := p.proj.GetConfig().Marshal()
	if err != nil {
		return errors.Wrap(err, "error building config")
	}

	err = os.WriteFile(p.proj.GetConfig().Path, b, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing config to file")
	}

	err = go_actions.PrepareGoConfigFolderAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	err = go_actions.TidyAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error building golang config")
	}

	err = git.ForceCommit(p.proj.GetProjectPath(), "added "+strings.Join(args, "; "))
	if err != nil {
		return errors.Wrap(err, "error performing git commit")
	}

	return nil
}

func (p *projectAdd) loadProject(cmd *cobra.Command) (err error) {
	pathToProject := cmd.Flag(pathFlag).Value.String()
	if pathToProject == "" {
		pathToProject = p.path
	}

	p.proj, err = project.LoadProject(pathToProject, p.config)
	if err != nil {
		return err
	}

	return nil
}
