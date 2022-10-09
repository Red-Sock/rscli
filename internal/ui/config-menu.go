package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/input"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli-uikit/multiselect"
	"github.com/Red-Sock/rscli-uikit/selectone"
	"github.com/Red-Sock/rscli/internal/randomizer"
	config2 "github.com/Red-Sock/rscli/pkg/service/config"
	"os"
)

const (
	pgCon    = "pg connection"
	redisCon = "redis connection"
)

func newConfigMenu() uikit.UIElement {
	msb := multiselect.New(
		configCallback,
		multiselect.Items(pgCon, redisCon),
	)

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

	cfg, err := config2.Run(args)
	if err != nil {
		return label.New(err.Error())
	}

	confDiag := cfgDialog{cfg: cfg}

	tb := input.New(confDiag.selectPathForConfig)
	// TODO RSI-23
	tb.W = 10
	tb.H = 1
	return tb
}

type cfgDialog struct {
	cfg *config2.Config
}

func (c *cfgDialog) selectPathForConfig(p string) uikit.UIElement {
	if p != "" {
		err := c.cfg.SetPath(p)
		if err != nil {
			if err == os.ErrExist {
				sb := selectone.New(
					c.processOverrideAnswer,
					selectone.Items("yes", "no"),
					selectone.Header("file "+p+" already exists. Want to override?"),
				)
				return sb
			}
		}
	}

	err := c.cfg.TryWrite()
	if err != nil {
		if err == os.ErrExist {
			sb := selectone.New(
				c.processOverrideAnswer,
				selectone.Items("yes", "no"),
				selectone.Header("file "+c.cfg.GetPath()+" already exists. Want to override?"),
			)
			return sb
		}
	}
	return label.New("successfully created file at " + c.cfg.GetPath() + ". " + randomizer.GoodGoodBuy())
}

func (c *cfgDialog) processOverrideAnswer(answ string) uikit.UIElement {
	if answ == "yes" {
		err := c.cfg.ForceWrite()
		if err != nil {
			return label.New(err.Error())
		}
		return label.New("successfully created file at " + c.cfg.GetPath() + ". " + randomizer.GoodGoodBuy())
	}
	return label.New("aborting config creation. " + randomizer.GoodGoodBuy())
}
