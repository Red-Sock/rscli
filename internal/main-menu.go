package internal

import (
	"fmt"
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

func RunUI(args map[string][]string) error {

	if len(args) > 1 {
		return fmt.Errorf("too manu arguments. Need 0 or 1, got %d", len(args))
	}

	mm := mainMenu()

	var startScreen uikit.UIElement

	for item := range args {
		startScreen = pluginsWithUI[item].Run(mm)
	}

	if startScreen == nil {
		startScreen = mm
	}

	qE := make(chan struct{})

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)

		<-sig

		qE <- struct{}{}
	}()

	uikit.NewHandler(startScreen).Start(qE)

	return nil
}

func mainMenu() uikit.UIElement {
	items := make([]string, 0, len(pluginsWithUI)+1)

	for name := range pluginsWithUI {
		items = append(items, name)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
	items = append(items, "Exit")

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
	if resp == "Exit" {
		return nil
	}
	return pluginsWithUI[resp].Run(mainMenu())
}
