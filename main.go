package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment"
	initCmd "github.com/Red-Sock/rscli/cmd/project"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
)

func main() {
	root := &cobra.Command{
		Use: "rscli [command] [arguments] [flags]",

		Short: "RsCLI is a tool for handling developers environment",

		Version: "0.0.28-alpha",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: config.InitConfig,
		SilenceErrors:     true,
		SilenceUsage:      true,
	}

	root.PersistentFlags().String(config.CustomPathToConfig, "", "path flag to custom config")

	root.AddCommand(initCmd.NewCmd())
	root.AddCommand(environment.NewCmd())

	if err := root.Execute(); err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))
	}
}
