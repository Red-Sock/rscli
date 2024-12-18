package environment

import (
	"path"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/environment"
)

func newTidyEnvCmd(io io.IO, cfg *config.RsCliConfig) *cobra.Command {
	et := &envTidy{
		io:  io,
		cfg: cfg,
	}
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		RunE: et.RunTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(environment.PathFlag, environment.PathFlag[:1], "", `Path to folder with projects`)
	c.Flags().BoolP(environment.ServiceInContainer, environment.ServiceInContainer[:1], false, "Service will be run in container")

	return c
}

type envTidy struct {
	io  io.IO
	cfg *config.RsCliConfig
}

func (e *envTidy) RunTidy(cmd *cobra.Command, arg []string) error {
	e.io.Println("Running rscli env tidy")

	constructor, err := environment.NewGlobalEnv(e.io, e.cfg, e.getEnvDirPath(cmd))
	if err != nil {
		return rerrors.Wrap(err, "error creating global environment struct")
	}
	if !constructor.IsEnvExist() {
		err := constructor.Init()
		if err != nil {
			return rerrors.Wrap(err, "error initializing environment")
		}
	}

	return constructor.Tidy()
}

func (e *envTidy) getEnvDirPath(cmd *cobra.Command) string {
	envDirPath := cmd.Flag(environment.PathFlag).Value.String()

	if envDirPath == "" {
		envDirPath = io.GetWd()
	}

	if path.Base(envDirPath) != envpatterns.EnvDir {
		envDirPath = path.Join(envDirPath, envpatterns.EnvDir)
	}

	return envDirPath
}
