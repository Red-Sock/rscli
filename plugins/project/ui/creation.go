package ui

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	config "github.com/Red-Sock/rscli/plugins/config/pkg/const"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"

	managerConfig "github.com/Red-Sock/rscli/plugins/config"
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
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Want to create config or use existing?")),
		radioselect.Items(createConfig, useExistingConfig),
	)
}
func (c *createArgs) callbackForConfigSelect(resp string) uikit.UIElement {
	confirmCreation := radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Confirm creating project")),
		radioselect.Items(yes, noo),
	)

	switch resp {
	case createConfig:
		dir, _ := os.Getwd()

		var err error
		c.p.CfgPath = path.Join(dir, config.FileName)
		if err != nil {
			return shared_ui.GetHeaderFromText(err.Error())
		}

		return managerConfig.Run(confirmCreation)
	case useExistingConfig:
		return c.handleExistingConfig()

	default:
		return confirmCreation
	}
}

func (c *createArgs) handleExistingConfig() uikit.UIElement {
	dir, err := os.Getwd()
	if err != nil {
		return shared_ui.GetHeaderFromLabel(
			label.New(
				err.Error(),
				label.NextScreen(
					radioselect.New(
						c.callbackForConfigSelect,
						radioselect.HeaderLabel(shared_ui.GetHeaderFromText(fmt.Sprintf("Want to create config or use existing?"))),
						radioselect.Items(createConfig, useExistingConfig)))))
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return shared_ui.GetHeaderFromLabel(label.New(
			err.Error(),
			label.NextScreen(
				radioselect.New(
					c.callbackForConfigSelect,
					radioselect.HeaderLabel(shared_ui.GetHeaderFromText(fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name))),
					radioselect.Items(createConfig, useExistingConfig),
				))))
	}

	potentialConfigs := make([]string, 0, len(files)/2)

	for _, item := range files {
		name := item.Name()
		if strings.HasSuffix(name, yamlExtension) {
			potentialConfigs = append(potentialConfigs, path.Join(dir, name))
		}
	}

	sort.Slice(potentialConfigs, func(i, j int) bool {
		return potentialConfigs[i] > potentialConfigs[j]
	})

	return radioselect.New(
		c.callbackExistingConfig,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Select one of the files:")),
		radioselect.Items(potentialConfigs...),
	)
}
func (c *createArgs) callbackExistingConfig(answ string) uikit.UIElement {
	if answ == "" {
		return radioselect.New(
			c.callbackForConfigSelect,
			radioselect.HeaderLabel(shared_ui.GetHeaderFromText(fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name))),
			radioselect.Items(createConfig, useExistingConfig),
		)
	}
	var err error
	c.p.CfgPath = answ
	if err != nil {
		return shared_ui.GetHeaderFromText(err.Error())
	}

	return radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Confirm creating project")),
		radioselect.Items(yes, noo),
	)
}

func (c *createArgs) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == yes {
		proj, err := processor.CreateProject(c.p)
		if err != nil {
			return shared_ui.GetHeaderFromText(err.Error())
		}

		err = proj.Build()
		if err != nil {
			return shared_ui.GetHeaderFromText(err.Error())
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
