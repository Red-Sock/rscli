package ui

import (
	"fmt"
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/input"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/pkg/service/project"
	"os"
	"path"
)

const (
	projCreate = "create"

	// TODO

	projUpdate = "update" // update version
	projAdd    = "add"    // add new (source type, transport, etc)
)

func newProjectMenu() uikit.UIElement {
	sb := radioselect.New(
		projectCallback,
		radioselect.Header(help.Header+"Creating project"),
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
		return input.New(
			p.callBackInputName,
			input.Width(20),
			input.Height(1),
			input.Position(common.NewRelativePositioning(0.5, 0.5)),
			input.TextAbove(err.Error()+". Input project name"),
			input.TextBelow("Enter to confirm"),
		)
	}

	if p.p.CfgPath == "" {
		return radioselect.New(
			p.doConfig,
			radioselect.Header(help.Header+"Want to create config?"),
			radioselect.Items("yes", "no"),
		)
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

func (p *projectInteraction) doConfig(resp string) uikit.UIElement {
	confirmCreation := radioselect.New(
		p.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		radioselect.Items("yes", "no"),
	)
	if resp == "yes" {
		dir, _ := os.Getwd()
		p.p.CfgPath = path.Join(dir, config.FileName)

		return newConfigMenu(confirmCreation)
	}
	return confirmCreation
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
