package tg

import (
	"context"

	"github.com/Red-Sock/go_tg/client"

	"financial-microservice/internal/config"
	"financial-microservice/internal/transport/tg/handlers/version"
	"financial-microservice/internal/transport/tg/menus/mainmenu"
)

type Server struct {
	bot *client.Bot
}

func NewServer(cfg *config.Config) (s *Server) {
	s = &Server{}
	s.bot = client.NewBot(cfg.GetString(config.ServerTgApikey))

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg), version.Command)

	}

	{
		// Add pre-rendered  menus
		s.bot.AddMenu(mainmenu.NewMainMenu())
	}

	return s
}

func (s *Server) Start(_ context.Context) error {
	s.bot.Start()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.bot.Stop()
	return nil
}
