package project

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies/link_service"
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

	c.Flags().StringP(processor.PathFlag, processor.PathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectLink) run(_ *cobra.Command, args []string) (err error) {
	if p.proj == nil {
		p.proj, err = project.LoadProject(p.path, p.config)
		if err != nil {
			return errors.Wrap(err, "error fetching project for linking")
		}
	}

	p.io.Println("Linking project...")
	grpcClient := link_service.GrpcClient{
		Modules: args,
		Cfg:     p.config,
		Io:      p.io,
	}

	err = grpcClient.AppendToProject(p.proj)
	if err != nil {
		return errors.Wrap(err, "error applying grpc clients")
	}

	actionPerformer := actions.NewActionPerformer(p.io)

	err = actionPerformer.Tidy(p.proj)
	if err != nil {
		return errors.Wrap(err, "error tiding project")
	}

	p.io.Println("Tidy executed. Commiting changes")

	err = git.CommitWithUntracked(p.proj.GetProjectPath(), "added "+strings.Join(args, "; "))
	if err != nil {
		return errors.Wrap(err, "error performing git commit")
	}

	return nil
}
