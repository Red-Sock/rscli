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
	envConstr := newEnvConstructor()
	cmd.AddCommand(newInitEnvCmd(envConstr.runInit))
	cmd.AddCommand(newTidyEnvCmd(envConstr.runTidy))

	return cmd
}
