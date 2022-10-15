package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/input"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"
	"github.com/Red-Sock/rscli/pkg/service/project"
)

const (
	projCreate = "create"
)

func newProjectMenu() uikit.UIElement {
	sb := radioselect.New(
		projectCallback,
		radioselect.Items(projCreate),
	)

	return sb
}

type projectInteraction struct {
	p project.Project
}

func projectCallback(resp string) uikit.UIElement {
	var err error
	switch resp {
	case projCreate:
		pi := &projectInteraction{}

		pi.p, err = project.NewProject(nil)
		if err != nil {
			if err != project.ErrNoConfigNoAppNameFlag {
				return label.New(err.Error())
			}
			return projectNameTextBox(pi)
		}

		if pi.p.Name == "" {
			return projectNameTextBox(pi)
		}

		return radioselect.New(
			pi.confirmCreateProjectCallback,
			radioselect.Header(fmt.Sprintf("You wish to create project named %s", pi.p.Name)),
			radioselect.Items("yes", "no"),
		)
	}

	return nil
}

func (p *projectInteraction) callBackInputName(resp string) uikit.UIElement {
	p.p.Name = resp

	err := p.p.ValidateName()
	if err != nil {
		return projectNameTextBox(p)
	}

	so := radioselect.New(
		p.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		radioselect.Items("yes", "no"),
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

func projectNameTextBox(pi *projectInteraction) *input.TextBox {
	return input.New(
		pi.callBackInputName,
		input.Width(20),
		input.Height(1),
		input.Position(common.NewRelativePositioning(0.5, 0.5)),
		input.TextAbove("Input project name"),
		input.TextBelow("Enter to confirm"),
	)
}
