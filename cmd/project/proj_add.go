package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
)

const (
	postgresArgument = "postgres"
	redisArgument    = "redis"

	telegramArgument = "telegram"

	restArgument = "rest"
)

type dependency interface {
	Do(proj interfaces.Project) error
}

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

	deps := p.getDependenciesFromUser(args)
	for _, d := range deps {
		err = d.Do(p.proj)
		if err != nil {
			return errors.Wrap(err, "error adding dependency to project")
		}
	}

	err = p.proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building folders")
	}

	err = p.proj.GetConfig().BuildTo(p.proj.GetConfigPath())
	if err != nil {
		return errors.Wrap(err, "error building config")
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

func (p *projectAdd) getDependenciesFromUser(args []string) []dependency {
	serverOpts := make([]dependency, 0, len(args))

	for _, arg := range args {
		var dep dependency
		switch arg {
		case postgresArgument:
			dep = dependencies.Postgres{Cfg: p.config}
		case redisArgument:
			dep = dependencies.Redis{Cfg: p.config}
		case telegramArgument:
			dep = dependencies.Telegram{Cfg: p.config}
		case restArgument:
			dep = dependencies.Rest{Cfg: p.config}
		default:
			p.io.PrintlnColored(colors.ColorRed, "unknown dependency: "+arg)
		}

		serverOpts = append(serverOpts, dep)
	}

	return serverOpts
}
