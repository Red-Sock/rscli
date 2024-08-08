package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type Telegram struct {
	Name string

	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (t Telegram) GetFolderName() string {
	if t.Name != "" {
		return t.Name
	}

	return projpatterns.TelegramServer
}

func (t Telegram) AppendToProject(proj proj_interfaces.Project) error {
	err := t.applyClient(proj)
	if err != nil {
		return errors.Wrap(err, "error applying tg client")
	}

	err = t.applyFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying tg folder")
	}

	t.applyConfig(proj)

	return nil
}

func (t Telegram) applyClient(proj proj_interfaces.Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathsToClients, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	tgConnFile := projpatterns.TgConnFile.CopyWithNewName(
		path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), projpatterns.TgConnFile.Name))

	renamer.ReplaceProjectName(proj.GetName(), tgConnFile)

	proj.GetFolder().Add(
		tgConnFile,
	)

	return nil
}

func (t Telegram) applyFolder(proj proj_interfaces.Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathToServers, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	tgServer := projpatterns.TgServFile.Copy()

	renamer.ReplaceProjectName(proj.GetName(), tgServer)

	tgHandlerExample := projpatterns.TgHandlerExampleFile.Copy()

	renamer.ReplaceProjectName(proj.GetName(), tgHandlerExample)

	proj.GetFolder().Add(
		&folder.Folder{
			Name: path.Join(t.Cfg.Env.PathToServers[0], t.GetFolderName()),
			Inner: []*folder.Folder{
				tgServer,
				{
					Name:  path.Join(projpatterns.HandlersFolderName, projpatterns.VersionFolderName),
					Inner: []*folder.Folder{tgHandlerExample},
				},
			},
		},
	)

	return nil
}

func (t Telegram) applyConfig(proj proj_interfaces.Project) {
	for _, srv := range proj.GetConfig().DataSources {
		if srv.GetName() == t.GetFolderName() {
			return
		}
	}

	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, &resources.Telegram{
		Name: resources.Name(t.GetFolderName()),
	})
}
