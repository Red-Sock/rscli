package internal

import (
	"github.com/Red-Sock/rscli/internal/commands"
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/service/help"
)

type Tool interface {
	Run(args []string) (output string)
	HelpMessage() string
}

func init() {
	cfgTool := config.NewConfigTool()
	cmd[commands.Config] = cfgTool
	cmd["c"] = cfgTool
	cmd["cfg"] = cfgTool

	helpTool := help.NewHelpTool()
	cmd[commands.Help] = helpTool
}

var cmd = map[string]Tool{}

func Run(args []string) {
	if len(args) == 0 {
		println(cmd[commands.Help].Run(args))
		return
	}

	println(cmd[args[0]].Run(args))
}
