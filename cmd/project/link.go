package project

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
)

type projectLink struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
	path string
}

func newLinkCmd(pl projectLink) *cobra.Command {
	c := &cobra.Command{
		Use:   "link",
		Short: "Links other projects",
		Long:  `Can be used to link another project's contracts`,

		RunE: pl.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectLink) run(_ *cobra.Command, args []string) (err error) {
	if p.proj == nil {
		p.proj, err = project.LoadProject(p.path, p.config)
		if err != nil {
			return errors.Wrap(err, "error fetching project")
		}
	}

	err = dependencies.GrpcClient{
		Modules: args,
		Cfg:     p.config,
		Io:      p.io,
	}.AppendToProject(p.proj)
	if err != nil {
		return errors.Wrap(err, "error applying grpc clients")
	}

	err = go_actions.PrepareGoConfigFolderAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	err = go_actions.TidyAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error tiding project")
	}

	return nil
}
