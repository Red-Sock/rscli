package dependencies

import (
	"github.com/godverv/matreshka/server"
)

const defaultServerPort = 80

func prepareServer(proj Project) *server.Server {
	for _, s := range proj.GetConfig().Servers {
		return s
	}

	newServer := &server.Server{
		Name: "",
		GRPC: make(map[string]*server.GRPC),
		FS:   make(map[string]*server.FS),
		HTTP: make(map[string]*server.HTTP),
	}
	proj.GetConfig().Servers[defaultServerPort] = newServer

	return newServer
}
