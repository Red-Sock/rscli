package manager

import (
	"github.com/Red-Sock/rscli-uikit/composit-items/input"
	"os"
	"path"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/plugins/config/processor"
)

func Run(elem uikit.UIElement) uikit.UIElement {
	c := &cfgDialog{
		previousScreen: elem,
		flags:          map[string][]string{},
	}

	c.subMenus = map[string]*ConfigMenuSubItem{
		TransportTypeMenu: NewConfigMenuSubItem(TransportTypeItems(), c.mainMenu),
		DataSourceMenu:    NewConfigMenuSubItem(DataSourcesItems(), c.mainMenu),
	}

	return c.mainMenu()
}

type cfgDialog struct {
	cfg  *processor.Config
	path string

	previousScreen uikit.UIElement

	subMenus map[string]*ConfigMenuSubItem
	flags    map[string][]string
}

// main screen of config menu
func (c *cfgDialog) mainMenu() uikit.UIElement {
	return radioselect.New(
		c.mainMenuCallback,
		radioselect.Header("DataSources"),
		radioselect.Items(MainMenuItems()...),
	)
}

func (c *cfgDialog) mainMenuCallback(res string) uikit.UIElement {
	if res == CommitConfig {
		return c.askName()
	}

	subMenu, ok := c.subMenus[res]
	if !ok {
		return label.New("something went wrong 0_o")
	}

	return subMenu.UiElement()
}

func (c *cfgDialog) askName() uikit.UIElement {
	return input.New(
		c.nameCallback,
		input.Expandable(),
		input.TextAbove("Application name:"),
	)
}

func (c *cfgDialog) nameCallback(s string) uikit.UIElement {
	if s == "" {
		return c.askName()
	}

	c.flags["-"+processor.AppName] = []string{s}

	return c.commitConfig()
}

func (c *cfgDialog) commitConfig() uikit.UIElement {
	args := make([]string, 0, len(c.subMenus))
	for _, a := range c.subMenus {
		args = append(args, a.BuildFlagsForConfig()...)
	}

	for f, arg := range c.flags {
		args = append(args, f)
		args = append(args, arg...)
	}

	cfg, err := processor.Run(flag.ParseArgs(args))
	if err != nil {
		return label.New("error creating config: " + err.Error())
	}
	c.cfg = cfg

	c.path, _ = os.Getwd()
	c.path = path.Join(c.path, processor.FileName)

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
	return label.New("aborting config creation. ")
}

func (c *cfgDialog) endDialog() uikit.UIElement {
	if c.previousScreen == nil {
		return label.New("successfully created file at " + c.cfg.GetPath() + ". ")
	}

	return label.New("successfully created file at "+c.cfg.GetPath()+". (press enter to continue)",
		label.NextScreen(func() uikit.UIElement {
			return c.previousScreen
		}))
}
