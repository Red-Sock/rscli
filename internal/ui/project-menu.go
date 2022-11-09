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
	"sort"
	"strings"
)

const (
	projCreate = "create"

	// TODO

	projUpdate = "update" // update version
	projAdd    = "add"    // add new (source type, transport, etc)

	// config

	uec = "use existing config"
	yes = "yes"
	noo = "no"

	yamlExtension = ".yaml"
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
			radioselect.Items(yes, noo),
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
			radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", p.p.Name)),
			radioselect.Items(yes, noo, uec),
		)
	}

	so := radioselect.New(
		p.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		radioselect.Items(yes, noo),
	)
	return so
}

func (p *projectInteraction) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == yes {
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
		radioselect.Items(yes, noo),
	)
	switch resp {
	case yes:
		dir, _ := os.Getwd()
		p.p.CfgPath = path.Join(dir, config.FileName)

		return newConfigMenu(confirmCreation)
	case uec:
		return p.handleUEC()
	default:
		return confirmCreation
	}
}

func (p *projectInteraction) selectExistingConfig(answ string) uikit.UIElement {
	if answ == "" {
		return radioselect.New(
			p.doConfig,
			radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", p.p.Name)),
			radioselect.Items(yes, noo, uec),
		)
	}
	p.p.CfgPath = answ
	return radioselect.New(
		p.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", p.p.Name)),
		radioselect.Items(yes, noo),
	)
}

func (p *projectInteraction) handleUEC() uikit.UIElement {
	dir, err := os.Getwd()
	if err != nil {
		return label.New(
			err.Error(),
			label.NextScreen(func() uikit.UIElement {
				return radioselect.New(
					p.doConfig,
					radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", p.p.Name)),
					radioselect.Items(yes, noo, uec),
				)
			}))
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return label.New(
			err.Error(),
			label.NextScreen(func() uikit.UIElement {
				return radioselect.New(
					p.doConfig,
					radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", p.p.Name)),
					radioselect.Items(yes, noo, uec),
				)
			}))
	}

	potentialConfigs := make([]string, 0, len(files)/2)
	otherFiles := make([]string, 0, len(files)/2)

	for _, item := range files {
		name := item.Name()
		if strings.HasPrefix(name, yamlExtension) {
			potentialConfigs = append(potentialConfigs, path.Join(dir, name))
		} else {
			otherFiles = append(otherFiles, path.Join(dir, name))
		}
	}

	sort.Slice(potentialConfigs, func(i, j int) bool {
		return potentialConfigs[i] > potentialConfigs[j]
	})

	sort.Slice(otherFiles, func(i, j int) bool {
		return otherFiles[i] > otherFiles[j]
	})

	return radioselect.New(
		p.selectExistingConfig,
		radioselect.Header("Select one of the files:"),
		radioselect.Items(append(potentialConfigs, otherFiles...)...),
	)
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
