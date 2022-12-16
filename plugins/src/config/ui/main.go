package main

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/plugins/src/config/ui/manager"
)

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return "config"
}

func (p *plugin) Run(element uikit.UIElement) uikit.UIElement {
	return manager.Run(element)
}
