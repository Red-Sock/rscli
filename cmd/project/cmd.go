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

	cmd.AddCommand(newInitCmd(projectInit{
		io:     io.StdIO{},
		config: config.GetConfig(),
	}))

	cmd.AddCommand(newAddCmd())

	return cmd
}
