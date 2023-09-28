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
	defaultApiName := r.GetFolderName() + "_api"

	err := r.applyFolder(proj, defaultApiName)
	if err != nil {
		return errors.Wrap(err, "error applying rest folder")
	}

	r.applyConfig(proj, defaultApiName)

	return nil
}

func (r Rest) applyConfig(proj interfaces.Project, defaultApiName string) {
	ds := proj.GetConfig().Server

	if _, ok := ds[defaultApiName]; ok {
		return
	}

	ds[defaultApiName] = server.Rest{}

}

func (r Rest) applyFolder(proj interfaces.Project, defaultApiName string) error {
	ok, err := containsDependency(r.Cfg.Env.PathToServers, proj.GetFolder(), r.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error searching dependencies")
	}

	if ok {
		return nil
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, patterns.ServerGoFile),
			Content: patterns.RestServFile,
		},
		&folder.Folder{
			Name:    path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, patterns.VersionGoFile),
			Content: patterns.RestServHandlerExampleFile,
		},
	)

	return nil
}