package internal

import (
	"os"
	"os/signal"
	"sort"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/endscreen"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/service/help"

	cfgui "github.com/Red-Sock/rscli/plugins/config"
	envui "github.com/Red-Sock/rscli/plugins/environment"
	projectui "github.com/Red-Sock/rscli/plugins/project"
)

var plugins = map[string]func(uikit.UIElement) uikit.UIElement{
	cfgui.PluginName:     cfgui.Run,
	projectui.PluginName: projectui.Run,
	envui.PluginName:     envui.Run,
}

func Run() error {
	qE := make(chan struct{})

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)

		<-sig

		qE <- struct{}{}
	}()

	uikit.NewHandler(mainMenu()).Start(qE)

	return nil
}

func mainMenu() uikit.UIElement {
	items := make([]string, 0, len(plugins)+1)

	for name := range plugins {
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
		getMainMenuCallback,
		radioselect.Header(help.Header+"Main menu"),
		radioselect.Items(items...),
		radioselect.PreviousScreen(&endscreen.EndScreen{UIElement: label.New(randomizer.GoodGoodBuy())}),
	)
}

func getMainMenuCallback(resp string) uikit.UIElement {
	switch resp {
	case cfgui.PluginName:
		return cfgui.Run(mainMenu())
	case projectui.PluginName:
		return projectui.Run(mainMenu())
	case envui.PluginName:
		return envui.Run(mainMenu())
	case "exit":
		return nil
	default:
		return nil
	}
}
