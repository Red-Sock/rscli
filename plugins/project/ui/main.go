package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/plugins/project/ui/manager"
)

const PluginName = "project"

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return PluginName
}

func (p *plugin) Run(element uikit.UIElement) uikit.UIElement {
	return manager.Run(element)
}
