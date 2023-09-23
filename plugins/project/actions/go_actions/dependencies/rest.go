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
		return ErrNoServerFolderInConfig
	}
	defaultApiName := r.GetFolderName() + "_api"

	r.doConfig(proj, defaultApiName)
	r.doFolder(proj, defaultApiName)

	err := proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	return nil
}

func (r Rest) doConfig(proj interfaces.Project, defaultApiName string) {
	ds := proj.GetConfig().Server
	if ds == nil {
		ds = make(map[string]interface{})
	}

	containsConfig := false
	for name := range ds {
		if name == defaultApiName {
			containsConfig = true
			break
		}
	}

	if containsConfig {
		return
	}
	ds[defaultApiName] = server.Rest{}

	proj.GetConfig().Server = ds
}

func (r Rest) doFolder(proj interfaces.Project, defaultApiName string) {
	containsFolder := false
	for _, pth := range r.Cfg.Env.PathToServers {
		if proj.GetFolder().GetByPath(pth, defaultApiName) != nil {
			containsFolder = true
			break
		}
	}

	if containsFolder {
		return
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

}
