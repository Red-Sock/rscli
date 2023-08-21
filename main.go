package main

import (
	"fmt"

	"github.com/spf13/cobra"

	initCmd "github.com/Red-Sock/rscli/cmd/init"
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/pkg/colors"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/go_actions/update"
)

func main() {
	root := &cobra.Command{
		Use: "rscli [command] [arguments] [flags]",

		Short: "RsCLI is a tool for handling developers environment",
		Long:  "RsCLI is a useful developer tool for project generation and developer environment handling",

		Version: update.GetLatestVersion().String(),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(initCmd.Cmd)

	if err := root.Execute(); err != nil {
		stdio.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))
	}
}
