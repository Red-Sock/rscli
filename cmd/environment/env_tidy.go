package environment

import (
	"github.com/spf13/cobra"
)

func newTidyEnvCmd(command func(cmd *cobra.Command, _ []string) error) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		RunE: command,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `Path to folder with projects`)

	return c
}

func (c *envConstructor) runTidy(cmd *cobra.Command, _ []string) error {
	// TODO
	return nil
}
