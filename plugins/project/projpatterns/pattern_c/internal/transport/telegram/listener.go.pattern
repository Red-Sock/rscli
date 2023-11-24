package telegram

import (
	"context"

	"github.com/Red-Sock/go_tg/client"

	"proj_name/internal/config"
	"proj_name/internal/transport/telegram/version"
)

type Server struct {
	bot *client.Bot
}

func NewServer(cfg config.Config, bot *client.Bot) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg))

	}

	return s
}

func (s *Server) Start(_ context.Context) error {
	return s.bot.Start()
}

func (s *Server) Stop(_ context.Context) error {
	s.bot.Stop()
	return nil
}
