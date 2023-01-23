package project

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/plugins/project/processor"
	"github.com/Red-Sock/rscli/plugins/project/scripts"
)

const PluginName = "project"

func Run(prev uikit.UIElement) uikit.UIElement {
	pm := projectMenu{
		previous: prev,
	}
	sb := radioselect.New(
		pm.selectAction,
		radioselect.Header(help.Header+"Creating project"),
		radioselect.Items(projCreate),
	)

	return sb
}

type projectMenu struct {
	p        *processor.Project
	previous uikit.UIElement
}

func (pm *projectMenu) selectAction(resp string) uikit.UIElement {
	switch resp {
	case projCreate:
		return scripts.StartCreateProj(pm.previous)
	}

	return nil
}

const (
	projCreate = "create"
)
