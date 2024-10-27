package telegram

import (
	"github.com/Red-Sock/go_tg"
	"github.com/godverv/matreshka/resources"
)

type Bot go_tg.Bot

func New(cfg *resources.Telegram) (*Bot, error) {
	bot := go_tg.NewBot(cfg.ApiKey)
	return (*Bot)(bot), nil
}
