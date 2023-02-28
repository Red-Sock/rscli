package ui

import "C"
import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/pkg/service/help"

	manager "github.com/Red-Sock/rscli/plugins/config"
	config "github.com/Red-Sock/rscli/plugins/config/processor"
	"github.com/Red-Sock/rscli/plugins/project/processor"
)

type createArgs struct {
	p processor.CreateArgs

	previousScreen uikit.UIElement
}

func StartCreateProj(ps uikit.UIElement) uikit.UIElement {
	ca := createArgs{
		previousScreen: ps,
	}

	return ca.configDiag()
}

func (c *createArgs) configDiag() uikit.UIElement {
	return radioselect.New(
		c.callbackForConfigSelect,
		radioselect.Header(help.Header+"Want to create config or use existing?"),
		radioselect.Items(createConfig, useExistingConfig),
	)
}
func (c *createArgs) callbackForConfigSelect(resp string) uikit.UIElement {
	confirmCreation := radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.Header("Confirm creating project"),
		radioselect.Items(yes, noo),
	)

	switch resp {
	case createConfig:
		dir, _ := os.Getwd()

		var err error
		c.p.CfgPath = path.Join(dir, config.FileName)
		if err != nil {
			return label.New(err.Error())
		}

		return manager.Run(confirmCreation)
	case useExistingConfig:
		return c.handleExistingConfig()

	default:
		return confirmCreation
	}
}

func (c *createArgs) handleExistingConfig() uikit.UIElement {
	dir, err := os.Getwd()
	if err != nil {
		return label.New(
			err.Error(),
			label.NextScreen(func() uikit.UIElement {
				return radioselect.New(
					c.callbackForConfigSelect,
					radioselect.Header(help.Header+fmt.Sprintf("Want to create config or use existing?")),
					radioselect.Items(createConfig, useExistingConfig),
				)
			}))
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return label.New(
			err.Error(),
			label.NextScreen(func() uikit.UIElement {
				return radioselect.New(
					c.callbackForConfigSelect,
					radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name)),
					radioselect.Items(createConfig, useExistingConfig),
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
		c.callbackExistingConfig,
		radioselect.Header("Select one of the files:"),
		radioselect.Items(append(potentialConfigs, otherFiles...)...),
	)
}
func (c *createArgs) callbackExistingConfig(answ string) uikit.UIElement {
	if answ == "" {
		return radioselect.New(
			c.callbackForConfigSelect,
			radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name)),
			radioselect.Items(createConfig, useExistingConfig),
		)
	}
	var err error
	c.p.CfgPath = answ
	if err != nil {
		return label.New(err.Error())
	}

	return radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.Header("Confirm creating project"),
		radioselect.Items(yes, noo),
	)
}

func (c *createArgs) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == yes {
		proj, err := processor.CreateProject(c.p)
		if err != nil {
			return label.New(err.Error())
		}

		err = proj.Build()
		if err != nil {
			return label.New(err.Error())
		}
	}
	return c.previousScreen
}

const (
	useExistingConfig = "use existing config"
	createConfig      = "create new config"
	yes               = "yes"
	noo               = "no"
	yamlExtension     = ".yaml"
)
