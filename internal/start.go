package internal

import (
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"sort"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/endscreen"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/internal/randomizer"
	cfgui "github.com/Red-Sock/rscli/plugins/config"
	envui "github.com/Red-Sock/rscli/plugins/environment"
	projectui "github.com/Red-Sock/rscli/plugins/project"
)

var ErrUnknownCommand = errors.New("unknown command")

var plugins = map[string]func(uikit.UIElement) uikit.UIElement{
	cfgui.PluginName:     cfgui.Run,
	projectui.PluginName: projectui.RunProjectCMD,
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

	if len(os.Args) == 1 {
		uikit.NewHandler(mainMenu()).Start(qE)
		return nil
	} else {
		hand, ok := h.handles[os.Args[1]]
		if !ok {
			return errors.Wrapf(ErrUnknownCommand, "%s is not a defined command", os.Args[1])
		}
		return hand.Do(os.Args[2:])
	}
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
		return shared_ui.GetHeaderFromText("no plugins available")
	}

	return radioselect.New(
		getMainMenuCallback,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Main menu")),
		radioselect.Items(items...),
		radioselect.PreviousScreen(
			&endscreen.EndScreen{
				UIElement: shared_ui.GetHeaderFromText(randomizer.GoodGoodBuy()),
			},
		),
	)
}

func getMainMenuCallback(resp string) uikit.UIElement {
	switch resp {
	case cfgui.PluginName:
		return cfgui.Run(mainMenu())
	case projectui.PluginName:
		return projectui.RunProjectCMD(mainMenu())
	case envui.PluginName:
		return envui.Run(mainMenu())
	case "exit":
		return nil
	default:
		return nil
	}
}
