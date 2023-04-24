package config

import (
	"os"
	"path"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/input"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"

	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/plugins/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/config/processor"
)

const PluginName = "config"

func Run(prevousMenu uikit.UIElement) uikit.UIElement {
	c := &cfgDialog{
		flags: map[string][]string{},
		prev:  prevousMenu,
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

	prev uikit.UIElement

	subMenus map[string]*ConfigMenuSubItem
	flags    map[string][]string
}

// main screen of config menu
func (c *cfgDialog) mainMenu() uikit.UIElement {
	return radioselect.New(
		c.mainMenuCallback,
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Project configuration")),
		radioselect.Items(MainMenuItems()...),
	)
}

func (c *cfgDialog) mainMenuCallback(res string) uikit.UIElement {
	if res == CommitConfig {
		return c.askName()
	}

	subMenu, ok := c.subMenus[res]
	if !ok {
		return shared_ui.GetHeaderFromText("something went wrong 0_o")
	}

	return subMenu.UiElement()
}

func (c *cfgDialog) askName() uikit.UIElement {
	return input.New(
		c.nameCallback,
		input.ExpandableWithMaxWidth(20),
		input.Position(common.NewRelativePositioning(0.4, 0.5)),
		input.TextAbove("Application name:"),
	)
}

func (c *cfgDialog) nameCallback(s string) uikit.UIElement {
	if s == "" {
		return c.askName()
	}

	c.flags["-"+_const.AppName] = []string{s}

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
		return shared_ui.GetHeaderFromText("error creating config: " + err.Error())
	}
	c.cfg = cfg

	c.path, _ = os.Getwd()
	c.path = path.Join(c.path, _const.FileName)

	err = c.cfg.SetFolderPath(c.path)
	if err != nil {
		if err == os.ErrExist {
			sb := radioselect.New(
				c.handleOverrideAnswer,
				radioselect.Items("yes", "no"),
				radioselect.HeaderLabel(
					shared_ui.GetHeaderFromText("file "+c.path+" already exists. Want to override?"),
				),
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
				radioselect.HeaderLabel(
					shared_ui.GetHeaderFromText("file "+c.cfg.GetPath()+" already exists. Want to override?")),
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
			return shared_ui.GetHeaderFromText(err.Error())
		}
		return c.endDialog()
	}
	return shared_ui.GetHeaderFromText("aborting config creation. ")
}

func (c *cfgDialog) endDialog() uikit.UIElement {
	return shared_ui.GetHeaderFromLabel(
		label.New("successfully created file at "+c.cfg.GetPath()+". (press enter to continue)",
			label.NextScreen(c.prev),
		))
}
