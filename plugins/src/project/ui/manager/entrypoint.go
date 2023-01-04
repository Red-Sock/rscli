package manager

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/plugins/src/project/processor"
)

func Run(element uikit.UIElement) uikit.UIElement {
	pm := projectMenu{
		previous: element,
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
		return StartCreateProj(pm.previous)
	}

	return nil
}

const (
	projCreate = "create"
)
