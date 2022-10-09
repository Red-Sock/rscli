package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/selectone"
	"github.com/Red-Sock/rscli/pkg/service/help"
)

const (
	configMenu  = "config"
	projectMenu = "project"
)

func newMainMenu() uikit.UIElement {
	sb := selectone.New(
		mainMenuCallback,
		selectone.Header(help.Header),
		selectone.Items(
			configMenu,
			projectMenu,
		),
	)
	return sb

}

func mainMenuCallback(res string) uikit.UIElement {
	switch res {
	case configMenu:
		return newConfigMenu()
	case projectMenu:
		return newProjectMenu()
	default:
		return nil
	}
}
