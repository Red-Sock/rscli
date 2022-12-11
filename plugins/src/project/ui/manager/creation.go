package manager

import "C"
import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/input"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/Red-Sock/rscli/pkg/service/project/config-processor/config"
	config_ui "github.com/Red-Sock/rscli/pkg/service/project/config-ui"
	processor "github.com/Red-Sock/rscli/pkg/service/project/project-processor"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/validators"
)

type createArgs struct {
	p processor.CreateArgs

	previousScreen uikit.UIElement
}

func StartCreateProj(ps uikit.UIElement) uikit.UIElement {
	ca := createArgs{
		previousScreen: ps,
	}

	return ca.textBoxForName()
}

func (c *createArgs) textBoxForName(prefix ...string) *input.TextBox {
	pref := strings.Join(append(prefix, ""), ".")

	return input.New(
		c.callBackForName,
		input.Width(20),
		input.Height(1),
		input.Position(common.NewRelativePositioning(0.5, 0.5)),
		input.TextAbove(pref+"Input project name"),
		input.TextBelow("Enter to confirm"),
	)
}

func (c *createArgs) callBackForName(resp string) uikit.UIElement {
	if resp == "" {
		return c.textBoxForName()
	}

	c.p.Name = resp

	err := validators.ValidateNameString(c.p.Name)
	if err != nil {
		return c.textBoxForName(err.Error())
	}

	return radioselect.New(
		c.callbackForConfigSelect,
		radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name)),
		radioselect.Items(yes, noo, useExistingConfig),
	)

}

func (c *createArgs) callbackForConfigSelect(resp string) uikit.UIElement {
	confirmCreation := radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", c.p.Name)),
		radioselect.Items(yes, noo),
	)

	switch resp {
	case yes:
		dir, _ := os.Getwd()

		var err error
		c.p.CfgPath = path.Join(dir, config.FileName)
		if err != nil {
			return label.New(err.Error())
		}

		return config_ui.NewConfigMenu(confirmCreation)
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
					radioselect.Header(help.Header+fmt.Sprintf("Want to create config to project named \"%s\"?", c.p.Name)),
					radioselect.Items(yes, noo, useExistingConfig),
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
					radioselect.Items(yes, noo, useExistingConfig),
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
			radioselect.Items(yes, noo, useExistingConfig),
		)
	}
	var err error
	c.p.CfgPath = answ
	if err != nil {
		return label.New(err.Error())
	}

	return radioselect.New(
		c.confirmCreateProjectCallback,
		radioselect.Header(fmt.Sprintf("You wish to create project named %s", c.p.Name)),
		radioselect.Items(yes, noo),
	)
}

func (c *createArgs) confirmCreateProjectCallback(resp string) uikit.UIElement {
	if resp == yes {
		proj, err := processor.New(c.p)
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
	yes               = "yes"
	noo               = "no"

	yamlExtension = ".yaml"
)
