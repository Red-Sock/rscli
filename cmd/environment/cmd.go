package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/environment/env"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Handles environment",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdIO := io.StdIO{}

	cfg := config.GetConfig()

	cmd.AddCommand(newInitEnvCmd(&envInit{
		io:          stdIO,
		constructor: env.NewConstructor(stdIO, cfg),
	}))

	cmd.AddCommand(newTidyEnvCmd(&envTidy{
		io:          stdIO,
		constructor: env.NewConstructor(stdIO, cfg),
	}))

	return cmd
}
