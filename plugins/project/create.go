package project

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	"github.com/Red-Sock/rscli/plugins/project/processor"
	"github.com/Red-Sock/rscli/plugins/project/ui"
)

const PluginName = "project"

func RunProjectCMD(prev uikit.UIElement) uikit.UIElement {
	pm := projectMenu{
		previous: prev,
	}
	sb := radioselect.New(
		pm.selectAction,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Creating project")),
		radioselect.Items(projCreate, projTidy),
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
		return ui.StartCreateProj(pm.previous)
	case projTidy:
		return Tidy(pm.previous)
	}

	return nil
}

const (
	projCreate = "create"
	projTidy   = "tidy"
)
