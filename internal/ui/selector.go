package ui

import (
	"fmt"
	"os"
	"os/signal"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/pkg/service/project"
)

const Command = "ui"

func Run(args []string) {
	if len(args) > 2 {
		help.FormMessage(helpUI)
	}

	qE := make(chan struct{})

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		<-sig
		qE <- struct{}{}
	}()

	uikit.NewHandler(selectMenu(slices.Exclude(args, Command))).Start(qE)
}

func selectMenu(args []string) uikit.UIElement {
	switch {
	case len(args) == 0:
		return newMainMenu()
	case slices.Contains(config.Command(), args[0]):
		return newConfigMenu(nil)
	case slices.Contains(project.Command(), args[0]):
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
