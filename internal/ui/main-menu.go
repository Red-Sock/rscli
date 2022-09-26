package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/selectone"
)

const Command = "ui"

type MainMenu struct {
	selectone.SelectBox
}

const (
	configMenu = "crete config"
)

func NewMainMenu() uikit.UIElement {
	sb, _ := selectone.New(
		mainMenuCallback,
		selectone.ItemsAttribute(configMenu),
	)
	return sb

}

func mainMenuCallback(res string) uikit.UIElement {
	switch res {
	case configMenu:
		return newConfigMenu()
	default:
		return nil
	}
}
