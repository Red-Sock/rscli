package internal

import (
	"os"
	"os/signal"

	"github.com/Red-Sock/rscli-uikit/basic/endscreen"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/service/help"

	uikit "github.com/Red-Sock/rscli-uikit"
)

func RunUI(args []string) {
	if len(args) == 0 {
		println("no arguments given")
		return
	}

	qE := make(chan struct{})

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		<-sig
		qE <- struct{}{}
	}()

	uikit.NewHandler(mainMenu()).Start(qE)
}

func mainMenu() uikit.UIElement {
	items := make([]string, 0, len(pluginsWithUI)+1)

	for name := range pluginsWithUI {
		items = append(items, name)
	}

	if len(items) == 0 {
		return label.New(help.Header + "no plugins available")
	}

	return radioselect.New(
		mainMenuCallback,
		radioselect.Header(help.Header+"Main menu"),
		radioselect.Items(items...),
		radioselect.PreviousScreen(&endscreen.EndScreen{UIElement: label.New(randomizer.GoodGoodBuy())}),
	)
}

func mainMenuCallback(resp string) uikit.UIElement {
	return pluginsWithUI[resp].Run(mainMenu())
}
