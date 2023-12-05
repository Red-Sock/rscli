package project

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
)

const (
	nameFlag = "name"
	pathFlag = "path"
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

	cmd.AddCommand(newInitCmd(projectInit{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	cmd.AddCommand(newAddCmd(projectAdd{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	cmd.AddCommand(newLinkCmd(projectLink{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	return cmd
}
