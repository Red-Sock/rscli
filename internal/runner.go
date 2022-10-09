package internal

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/internal/cmds"
	"github.com/Red-Sock/rscli/internal/ui"
	"github.com/Red-Sock/rscli/internal/utils"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/pkg/service/project"
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

	if utils.Contains(args, ui.Command) {
		qE := make(chan struct{})
		uikit.NewHandler(ui.NewUI(args)).Start(qE)
		return
	}

	switch {
	case utils.Contains(config.Command(), args[0]):
		cmds.RunConfig(args[1:])
	case utils.Contains(project.Command(), args[0]):
		cmds.RunProject(args[1:])
	}

	println(res)
}
