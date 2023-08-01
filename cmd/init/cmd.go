package init

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(projectCmd)
}

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes project",
	Long:  `Can be used to init project in RSCLI project style`,
}
