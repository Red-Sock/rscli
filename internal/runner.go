package internal

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/internal/cmds"
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/service/help"
	"github.com/Red-Sock/rscli/internal/ui"
	"github.com/Red-Sock/rscli/internal/utils"
)

type Tool interface {
	Run(args []string) (output string)
	HelpMessage() string
}

func Run(args []string) {
	if len(args) == 0 {
		println(help.Run())
		return
	}
	var res string
	switch {
	case utils.Contains(config.Command(), args[0]):
		cmds.RunConfig(args[0:])
	case utils.Contains(args, ui.Command):
		qE := make(chan struct{})
		uikit.NewHandler(ui.NewUI(args)).Start(qE)
		return
	}

	println(res)
}
