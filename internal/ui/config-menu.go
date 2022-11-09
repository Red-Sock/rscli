package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"os"
	"path"
)

// menu itself
func newConfigMenu(previousSequence uikit.UIElement) uikit.UIElement {
	c := &cfgDialog{
		previousScreen: previousSequence,
	}

	c.subMenus = map[string]*configMenuSubItem{
		transportTypeMenu: newConfigMenuSubItem(transportTypeItems(), c.configMenu),
		dataSourceMenu:    newConfigMenuSubItem(dataSourcesItems(), c.configMenu),
	}

	return c.configMenu()
}

type cfgDialog struct {
	cfg  *config.Config
	path string

	previousScreen uikit.UIElement

	subMenus map[string]*configMenuSubItem
}

// main screen of config menu
func (c *cfgDialog) configMenu() uikit.UIElement {
	return radioselect.New(
		c.mainMenuCallback,
		radioselect.Header(help.Header+"DataSources"),
		radioselect.Items(mainMenuItems()...),
	)
}

func (c *cfgDialog) mainMenuCallback(res string) uikit.UIElement {
	if res == commitConfig {
		return c.commitConfig()
	}

	subMenu, ok := c.subMenus[res]
	if !ok {
		return label.New("something went wrong 0_o")
	}

	return subMenu.uiElement()
}

func (c *cfgDialog) commitConfig() uikit.UIElement {
	args := make([]string, 0, len(c.subMenus))
	for _, a := range c.subMenus {
		args = append(args, a.buildFlagsForConfig()...)
	}

	cfg, err := config.Run(args)
	if err != nil {
		return label.New("error creating config: " + err.Error())
	}

	c.cfg = cfg

	c.path, _ = os.Getwd()
	c.path = path.Join(c.path, config.FileName)

	err = c.cfg.SetFolderPath(c.path)
	if err != nil {
		if err == os.ErrExist {
			sb := radioselect.New(
				c.handleOverrideAnswer,
				radioselect.Items("yes", "no"),
				radioselect.Header("file "+c.path+" already exists. Want to override?"),
			)
			return sb
		}
	}

	err = c.cfg.TryWrite()
	if err != nil {
		if err == os.ErrExist {
			sb := radioselect.New(
				c.handleOverrideAnswer,
				radioselect.Items("yes", "no"),
				radioselect.Header("file "+c.cfg.GetPath()+" already exists. Want to override?"),
			)
			return sb
		}
	}
	return c.endDialog()
}

func (c *cfgDialog) handleOverrideAnswer(answ string) uikit.UIElement {
	if answ == "yes" {
		err := c.cfg.ForceWrite()
		if err != nil {
			return label.New(err.Error())
		}
		return c.endDialog()
	}
	return label.New("aborting config creation. " + randomizer.GoodGoodBuy())
}

func (c *cfgDialog) endDialog() uikit.UIElement {
	if c.previousScreen == nil {
		return label.New("successfully created file at " + c.cfg.GetPath() + ". " + randomizer.GoodGoodBuy())
	}
	return label.New("successfully created file at "+c.cfg.GetPath()+". ",
		label.NextScreen(func() uikit.UIElement {
			return c.previousScreen
		}))
}
