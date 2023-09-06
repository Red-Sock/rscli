package environment

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Handles environment",

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(newInitEnvCmd())
	cmd.AddCommand(newTidyEnvCmd())

	return cmd
}
