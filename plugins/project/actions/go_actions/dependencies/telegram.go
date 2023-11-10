package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type Telegram struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (t Telegram) GetFolderName() string {
	return projpatterns.TelegramServer
}

func (t Telegram) Do(proj interfaces.Project) error {
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

func (t Telegram) applyClient(proj interfaces.Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathsToClients, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		return nil
	}

	tgConnFile := projpatterns.TgConnFile.CopyWithNewName(
		path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), projpatterns.TgConnFile.Name))

	go_actions.ReplaceProjectName(proj.GetName(), tgConnFile)

	proj.GetFolder().Add(
		tgConnFile,
	)

	return nil
}

func (t Telegram) applyFolder(proj interfaces.Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathToServers, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	tgServer := projpatterns.TgServFile.Copy()

	go_actions.ReplaceProjectName(proj.GetName(), tgServer)

	tgHandlerExample := projpatterns.TgHandlerExampleFile.Copy()

	go_actions.ReplaceProjectName(proj.GetName(), tgHandlerExample)

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

func (t Telegram) applyConfig(proj interfaces.Project) {
	for _, srv := range proj.GetConfig().Servers {
		if srv.GetName() == t.GetFolderName() {
			return
		}
	}

	proj.GetConfig().Resources = append(proj.GetConfig().Resources, &resources.Telegram{
		Name: resources.Name(t.GetFolderName()),
	})
}
