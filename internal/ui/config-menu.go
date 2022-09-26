package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/input"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli-uikit/multiselect"
	"github.com/Red-Sock/rscli-uikit/selectone"
	"github.com/Red-Sock/rscli/internal/service/config"
	"log"
	"os"
)

const (
	pgCon    = "pg connection"
	redisCon = "redis connection"
)

func newConfigMenu() uikit.UIElement {
	msb, err := multiselect.New(
		configCallback,
		multiselect.ItemsAttribute(pgCon, redisCon),
	)

	if err != nil {
		log.Fatal("error creating config selector", err)
	}
	return msb
}

func configCallback(res []string) uikit.UIElement {
	args := make([]string, 0, len(res))
	for _, item := range res {
		switch item {
		case pgCon:
			args = append(args, "--pg")
		case redisCon:
			args = append(args, "--rds")
		}
	}

	cfg, err := config.Run(args)
	if err != nil {
		return label.New(err.Error())
	}

	confDiag := cfgDialog{cfg: cfg}

	tb := input.NewTextBox(confDiag.selectPathForConfig)
	// TODO RSI-23
	tb.W = 10
	tb.H = 1
	return tb
}

type cfgDialog struct {
	cfg *config.Config
}

func (c *cfgDialog) selectPathForConfig(p string) uikit.UIElement {
	if p != "" {
		err := c.cfg.SetPath(p)
		if err != nil {
			if err == os.ErrExist {
				sb, _ := selectone.New(
					c.processOverrideAnswer,
					selectone.ItemsAttribute("yes", "no"),
					selectone.HeaderAttribute("file "+p+" already exists. Want to override?"),
				)
				return sb
			}
		}
	}

	err := c.cfg.TryWrite()
	if err != nil {
		if err == os.ErrExist {
			sb, _ := selectone.New(
				c.processOverrideAnswer,
				selectone.ItemsAttribute("yes", "no"),
				selectone.HeaderAttribute("file "+c.cfg.GetPath()+" already exists. Want to override?"),
			)
			return sb
		}
	}
	return label.New("successfully created file at " + c.cfg.GetPath())
}

func (c *cfgDialog) processOverrideAnswer(answ string) uikit.UIElement {
	if answ == "yes" {
		err := c.cfg.ForceWrite()
		if err != nil {
			return label.New(err.Error())
		}
		return label.New("successfully created file at " + c.cfg.GetPath())
	}
	return label.New("aborting config creation. see ya!")
}
