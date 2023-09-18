package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/env"
)

func newTidyEnvCmd(constr *env.Constructor) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		PreRunE: constr.FetchConstructor,
		RunE:    constr.RunTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(env.PathFlag, env.PathFlag[:1], "", `Path to folder with projects`)

	return c
}
