package project

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/io"
)

func newAddCmd() *cobra.Command {
	constr := projectAdd{
		io: io.StdIO{},
	}

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

type projectAdd struct {
	io io.IO
}

func (p *projectAdd) run(cmd *cobra.Command, args []string) error {
	p.getDependenciesFromUser(cmd, args)
	return nil
}

func (p *projectAdd) getDependenciesFromUser(cmd *cobra.Command, args []string) {
	p.io.Println(args...)
}
