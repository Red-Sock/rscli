package ui

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/label"
	"github.com/Red-Sock/rscli-uikit/multiselect"
	"github.com/Red-Sock/rscli/internal/service/config"
	"log"
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
	return label.New(config.Run(args))
}
