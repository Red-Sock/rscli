package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/env"
)

func newInitEnvCmd(constr *env.Constructor) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Init environment for projects in given folder",

		PreRunE: constr.FetchConstructor,
		RunE:    constr.RunInit,

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	c.Flags().StringP(env.PathFlag, env.PathFlag[:1], "", `Path to folder with projects`)
	return c
}
