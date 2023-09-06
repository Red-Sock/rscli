package project

import (
	"github.com/spf13/cobra"
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
	cmd.AddCommand(newInitProjectCmd(newProjectConstructor().run))

	return cmd
}
