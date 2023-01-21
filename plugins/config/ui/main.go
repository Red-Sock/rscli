package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/plugins/config/ui/manager"
)

const PluginName = "config"

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return PluginName
}

func (p *plugin) Run(element uikit.UIElement) uikit.UIElement {
	return manager.Run(element)
}
