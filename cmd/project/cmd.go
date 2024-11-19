package project

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/project/add"
	"github.com/Red-Sock/rscli/cmd/project/init_new"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/processor"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Handles project",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdIO := io.StdIO{}
	wd := io.GetWd()
	cfg := config.GetConfig()

	basicProc := processor.New()

	cmd.AddCommand(init_new.NewCommand(basicProc))
	cmd.AddCommand(add.NewCommand(basicProc))

	cmd.AddCommand(newLinkCmd(projectLink{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	cmd.AddCommand(newTidyCmd(projectTidy{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	return cmd
}
