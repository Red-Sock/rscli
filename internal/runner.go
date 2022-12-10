package internal

import (
	"github.com/Red-Sock/rscli/internal/cmds"
	"github.com/Red-Sock/rscli/internal/ui"
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/project"
)

type Tool interface {
	Run(args []string) (output string)
	HelpMessage() string
}

func Run(args []string) {
	if len(args) == 0 || slices.Contains(args, ui.Command) {
		ui.Run(args)
		return
	}

	var res string

	switch {
	case slices.Contains(config.Command(), args[0]):
		cmds.RunConfig(args[1:])
	case slices.Contains(project.Command(), args[0]):
		cmds.RunProject(args[1:])
	}

	println(res)
}
