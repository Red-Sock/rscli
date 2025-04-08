package dependencies

import (
	"path"

	"go.vervstack.ru/matreshka/pkg/matreshka/server"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

const defaultServerPort = 80

func prepareServerConfig(proj Project) *server.Server {
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

func initServerManagerFiles(proj Project) {
	serverManagerPath := []string{
		patterns.InternalFolder,
		patterns.TransportFolder,
		patterns.ServerManager.Name,
	}

	if proj.GetFolder().GetByPath(serverManagerPath...) == nil {
		proj.GetFolder().Add(
			patterns.ServerManager.
				CopyWithNewName(path.Join(serverManagerPath...)))
	}
}
