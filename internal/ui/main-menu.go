package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/pkg/service/help"
)

const (
	configMenu  = "config"
	projectMenu = "project"
	helpMenu    = "help"
)

func newMainMenu() uikit.UIElement {
	sb := radioselect.New(
		mainMenuCallback,
		radioselect.Header(help.Header+"Main menu"),
		radioselect.Items(
			configMenu,
			projectMenu,
			helpMenu,
		),
	)
	return sb

}

func mainMenuCallback(res string) uikit.UIElement {
	switch res {
	case configMenu:
		return newConfigMenu(nil)
	case projectMenu:
		return newProjectMenu()
	default:
		return label.New(help.Run(), label.NextScreen(newMainMenu))
	}
}
