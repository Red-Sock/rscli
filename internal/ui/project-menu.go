package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/input"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli-uikit/selectone"
	"github.com/Red-Sock/rscli/internal/service/project"
)

const (
	projCreate = "create"
)

func newProjectMenu() uikit.UIElement {
	sb := selectone.New(
		projectCallback,
		selectone.Items(projCreate),
	)

	return sb
}

type projectInteraction struct {
	p *project.Project
}

func projectCallback(resp string) uikit.UIElement {
	switch resp {
	case projCreate:
		p, err := project.NewProject(nil)
		if err != nil {
			return label.New(err.Error())
		}

		pi := projectInteraction{p: p}

		if pi.p.Name == "" {
			return input.New(
				pi.callBackInputName,
				input.Width(20),
				input.Height(1),
			)
		}

		return selectone.New(
			pi.confirmCreateProjectCallback,
			selectone.Header(fmt.Sprintf("You wish to create project named %s", p.Name)),
			selectone.Items("yes", "no"),
		)
	}

	return nil
}

func (p *projectInteraction) callBackInputName(resp string) uikit.UIElement {
	p.p.Name = resp

	err := p.p.ValidateName()
	if err != nil {
		return input.New(
			p.callBackInputName,
			input.Width(20),
			input.Height(1),
		)
	}
	so := selectone.New(
		p.confirmCreateProjectCallback,
		selectone.Header(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		selectone.Items("yes", "no"),
	)
	return so
}

func (p *projectInteraction) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == "yes" {
		err := p.p.Create()
		if err != nil {
			return label.New(err.Error())
		}
	}
	return nil
}
