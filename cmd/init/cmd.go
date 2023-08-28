package init

import (
	"github.com/spf13/cobra"
)

const (
	nameFlag = "name"
	pathFlag = "path"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes project",
		Long:  `Can be used to init project in RSCLI project style`,

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(newInitProjectCmd(newProjectConstructor().runProjectInit))

	return cmd
}
