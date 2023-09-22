package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config/server"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type Rest struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (r Rest) GetFolderName() string {
	return "rest"
}

func (r Rest) Do(proj interfaces.Project) error {
	if len(r.Cfg.Env.PathToServers) == 0 {
		return ErrNoClientFolderInConfig
	}

	defaultApiName := r.GetFolderName() + "_api"
	_ = defaultApiName

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(r.Cfg.Env.PathsToClients[0], r.GetFolderName(), patterns.ConnFileName),
			Content: patterns.RestServFile,
		},
	)

	err := proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	ds := proj.GetConfig().Server

	containsConfig := false
	for _, v := range ds {
		_ = v
	}

	if containsConfig {
		ds[r.GetFolderName()] = server.Rest{}
	}

	return nil
}
