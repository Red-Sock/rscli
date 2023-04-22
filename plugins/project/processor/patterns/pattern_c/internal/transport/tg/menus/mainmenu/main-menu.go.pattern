package mainmenu

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/response/menu"

	"financial-microservice/internal/transport/tg/handlers/version"
)

const Command = "/start"

func NewMainMenu() interfaces.Menu {
	m := menu.NewSimple("Main menu", Command)

	m.AddButton("SayHello", version.Command)

	return m
}
