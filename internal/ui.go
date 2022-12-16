package internal

import (
	"github.com/Red-Sock/rscli/internal/commands"
	"os"
	"os/signal"
	"sort"

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

	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
	items = append(items, commands.Exit)

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
	if resp == commands.Exit {
		return nil
	}
	return pluginsWithUI[resp].Run(mainMenu())
}
