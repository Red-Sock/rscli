package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/service/help"
	"github.com/Red-Sock/rscli/internal/service/project"
	"github.com/Red-Sock/rscli/internal/utils"
)

const Command = "ui"

func NewUI(args []string) uikit.UIElement {
	if len(args) > 2 {
		help.FormMessage(helpUI)
	}

	args = utils.Exclude(args, Command)

	switch {
	case len(args) == 0:
		return newMainMenu()
	case utils.Contains(config.Command(), args[0]):
		return newConfigMenu()
	case utils.Contains(project.Command(), args[0]):
		return newProjectMenu()
	default:
		return label.New(fmt.Sprintf("no ui for %s is available", args[0]))
	}
}

const helpUI = `
Invalid parameters amount.
For calling UI execute rscli ui or rscli <service name> ui.`
