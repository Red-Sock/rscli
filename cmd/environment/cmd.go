package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/env"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Handles environment",

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(newInitEnvCmd(env.NewEnvConstructor()))
	cmd.AddCommand(newTidyEnvCmd(env.NewEnvConstructor()))

	return cmd
}
