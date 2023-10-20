package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Handles environment",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdIO := io.StdIO{}

	cmd.AddCommand(newTidyEnvCmd(stdIO, config.GetConfig()))

	return cmd
}
