package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"os"
	"path"
)

const (
	pgCon    = "pg connection"
	redisCon = "redis connection"
)

func newConfigMenu(previousSequence uikit.UIElement) uikit.UIElement {
	c := &cfgDialog{
		previousScreen: previousSequence,
	}

	msb := multiselect.New(
		c.dataSourcesSelectCallback,
		multiselect.Header(help.Header+"DataSources"),
		multiselect.Items(
			pgCon,
			redisCon,
		),
		multiselect.SeparatorChecked([]rune{'x'}),
	)

	return msb
}

type cfgDialog struct {
	cfg  *config.Config
	path string

	previousScreen uikit.UIElement
}

func (c *cfgDialog) dataSourcesSelectCallback(res []string) uikit.UIElement {
	args := make([]string, 0, len(res))
	for _, item := range res {
		switch item {
		case pgCon:
			args = append(args, "-"+config.SourceNamePg)
		case redisCon:
			args = append(args, "-"+config.SourceNameRds)
		}
	}

	cfg, err := config.Run(args)
	if err != nil {
		return label.New("error creating config: " + err.Error())
	}

	c.cfg = cfg

	c.path, _ = os.Getwd()
	c.path = path.Join(c.path, config.FileName)

	return c.trySelectPathForConfig()
}

func (c *cfgDialog) trySelectPathForConfig() uikit.UIElement {
	err := c.cfg.SetFolderPath(c.path)
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
