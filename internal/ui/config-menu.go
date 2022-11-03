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
	transportTypeMenu = "transport type"
	dataSourceMenu    = "data source"
	commitConfig      = "done"

	pgCon    = "pg connection"
	redisCon = "redis connection"

	restHttpType = "HTTP/rest"
	grpcType     = "grpc"
)

func newConfigMenu(previousSequence uikit.UIElement) uikit.UIElement {
	c := &cfgDialog{
		previousScreen: previousSequence,
		args:           make(map[string][]string),
	}
	return c.configMenu()
}

type cfgDialog struct {
	cfg  *config.Config
	path string

	previousScreen uikit.UIElement

	args map[string][]string
}

func (c *cfgDialog) configMenu() uikit.UIElement {
	return radioselect.New(
		c.selectWhatToConfig,
		radioselect.Header(help.Header+"DataSources"),
		radioselect.Items(
			transportTypeMenu,
			dataSourceMenu,
			commitConfig,
		),
	)
}

func (c *cfgDialog) selectWhatToConfig(res string) uikit.UIElement {
	switch res {
	case transportTypeMenu:
		return multiselect.New(
			c.transportLayerSelectCallback,
			multiselect.Header(help.Header+"DataSources"),
			multiselect.Items(
				restHttpType,
				grpcType),
			multiselect.SeparatorChecked([]rune{'x'}),
		)
	case dataSourceMenu:
		return multiselect.New(
			c.dataSourcesSelectCallback,
			multiselect.Header(help.Header+"DataSources"),
			multiselect.Items(
				pgCon,
				redisCon,
			),
			multiselect.SeparatorChecked([]rune{'x'}),
		)
	case commitConfig:
		args := make([]string, 0, len(c.args))
		for _, a := range c.args {
			args = append(args, a...)
		}
		cfg, err := config.Run(args)
		if err != nil {
			return label.New("error creating config: " + err.Error())
		}

		c.cfg = cfg

		c.path, _ = os.Getwd()
		c.path = path.Join(c.path, config.FileName)
		return c.trySelectPathForConfig()

	default:
		return label.New("something went wrong 0_o")
	}
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
	c.args[dataSourceMenu] = args

	return c.configMenu()
}

func (c *cfgDialog) transportLayerSelectCallback(res []string) uikit.UIElement {
	args := make([]string, 0, len(res))
	for _, item := range res {
		switch item {
		case restHttpType:
			args = append(args, "-"+config.RESTHTTPServer)
		case grpcType:
			args = append(args, "-"+config.GRPCServer)
		}
	}
	c.args[transportTypeMenu] = args

	return c.configMenu()
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
