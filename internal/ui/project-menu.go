package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/input"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli-uikit/selectone"
	"github.com/Red-Sock/rscli/internal/service/project"
	"log"
)

const (
	projCreate = "create"
)

func newProjectMenu() uikit.UIElement {
	sb, err := selectone.New(
		projectCallback,
		selectone.ItemsAttribute(projCreate),
	)

	if err != nil {
		log.Fatal("error creating config selector", err)
	}

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
			// header attribute in RSI-22
			inputer := input.NewTextBox(pi.callBackInputName)
			inputer.W = 20
			inputer.H = 1
			return inputer
		}

		so, _ := selectone.New(
			pi.confirmCreateProjectCallback,
			selectone.HeaderAttribute(fmt.Sprintf("You wish to create project named %s", p.Name)),
			selectone.ItemsAttribute("yes", "no"),
		)
		return so
	}

	return nil
}

func (p *projectInteraction) callBackInputName(resp string) uikit.UIElement {
	p.p.Name = resp

	err := p.p.ValidateName()
	if err != nil {
		// header attribute in RSI-22
		inputer := input.NewTextBox(p.callBackInputName)
		inputer.W = 20
		inputer.H = 1
		return inputer
	}
	so, _ := selectone.New(
		p.confirmCreateProjectCallback,
		selectone.HeaderAttribute(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		selectone.ItemsAttribute("yes", "no"),
	)
	return so
}

func (p *projectInteraction) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == "yes" {
		p.p.Create()
	}
	return nil
}
