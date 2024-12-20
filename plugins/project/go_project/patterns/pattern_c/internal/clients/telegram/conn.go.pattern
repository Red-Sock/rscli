package telegram

import (
	"github.com/Red-Sock/go_tg"
	"go.verv.tech/matreshka/resources"
)

func New(cfg *resources.Telegram) (*go_tg.Bot, error) {
	bot := go_tg.NewBot(cfg.ApiKey)
	return bot, nil
}
