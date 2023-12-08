package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/api"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	renamer2 "github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type Rest struct {
	Name string

	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (r Rest) GetFolderName() string {
	if r.Name != "" {
		return r.Name
	}

	return "rest"
}

func (r Rest) AppendToProject(proj interfaces.Project) error {
	err := r.applyFolder(proj, r.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error applying rest folder")
	}

	r.applyConfig(proj, r.GetFolderName())
	applyServerFolder(proj)
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

	serverF.Content = renamer.ReplaceProjectNameFull(serverF.Content, proj.GetName())

	renamer2.ReplaceProjectName(proj.GetName(), serverF)

	proj.GetFolder().Add(
		serverF,
		projpatterns.RestServHandlerVersionExampleFile.CopyWithNewName(
			path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, projpatterns.RestServHandlerVersionExampleFile.Name)),
	)

	return nil
}

func applyServerFolder(proj interfaces.Project) {
	serverManagerPath := []string{projpatterns.InternalFolder, projpatterns.TransportFolder, projpatterns.ServerManagerPatternFile.Name}
	if proj.GetFolder().GetByPath(serverManagerPath...) == nil {
		proj.GetFolder().Add(
			projpatterns.ServerManagerPatternFile.
				CopyWithNewName(path.Join(serverManagerPath...)))
	}
}
