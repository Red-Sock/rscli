package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/api"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
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

	for _, item := range proj.GetConfig().Servers {
		if item.GetName() == defaultApiName {
			return
		}
	}

	proj.GetConfig().Servers = append(proj.GetConfig().Servers,
		&api.Rest{
			Name: api.Name(defaultApiName),
			Port: api.DefaultRestPort,
		})
}

func (r Rest) applyFolder(proj interfaces.Project, defaultApiName string) error {
	ok, err := containsDependencyFolder(r.Cfg.Env.PathToServers, proj.GetFolder(), r.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error searching dependencies")
	}

	if ok {
		return nil
	}
	serverF := projpatterns.RestServFile.CopyWithNewName(
		path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, projpatterns.RestServFile.Name))

	serverF.Content = renamer.ReplaceProjectName(serverF.Content, proj.GetName())

	go_actions.ReplaceProjectName(proj.GetName(), serverF)

	proj.GetFolder().Add(
		serverF,
		projpatterns.RestServHandlerVersionExampleFile.CopyWithNewName(
			path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, projpatterns.RestServHandlerVersionExampleFile.Name)),
	)

	return nil
}
