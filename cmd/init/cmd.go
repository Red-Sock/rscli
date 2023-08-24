package init

import (
	"github.com/spf13/cobra"
)

const (
	nameFlag = "name"
	pathFlag = "path"
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes project",
	Long:  `Can be used to init project in RSCLI project style`,
}

func init() {
	Cmd.AddCommand(newProjectCmd())
}
