package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
)

const (
	redisArgument    = "redis"
	postgresArgument = "postgres"
)

type projectAdd struct {
	io     io.IO
	path   string
	config *rscliconfig.RsCliConfig

	proj *project.Project
}

func newAddCmd(constr projectAdd) *cobra.Command {

	c := &cobra.Command{
		Use:   "add",
		Short: "Adds dependency to project project",
		Long:  `Can be used to add a datasource or external API dependency to project`,

		RunE: constr.run,

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

	p.getDependenciesFromUser(cmd, args)
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

func (p *projectAdd) getDependenciesFromUser(_ *cobra.Command, args []string) {
	p.io.Println(args...)
}
