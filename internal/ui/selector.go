package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli/internal/utils"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/pkg/service/project"
	"os"
	"os/signal"
)

const Command = "ui"

func Run(args []string) {
	if len(args) > 2 {
		help.FormMessage(helpUI)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	qE := make(chan struct{})
	uikit.NewHandler(selectMenu(utils.Exclude(args, Command))).Start(qE)

	<-sig
	qE <- struct{}{}
}

func selectMenu(args []string) uikit.UIElement {
	switch {
	case len(args) == 0:
		return newMainMenu()
	case utils.Contains(config.Command(), args[0]):
		return newConfigMenu(nil)
	case utils.Contains(project.Command(), args[0]):
		return newProjectMenu()
	default:
		return label.New(fmt.Sprintf("no ui for %s is available", args[0]))
	}
}

const helpUI = `
Invalid parameters amount.
For calling UI execute 
rscli ui 
or 
rscli <service name> ui`
