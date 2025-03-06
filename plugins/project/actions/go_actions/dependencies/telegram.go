package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

type Telegram struct {
	dependencyBase
}

func telegram(dep dependencyBase) Dependency {
	return &Telegram{
		dep,
	}
}

func (t Telegram) GetFolderName() string {
	if t.Name != "" {
		return t.Name
	}

	return patterns.TelegramServer
}

func (t Telegram) AppendToProject(proj Project) error {
	err := t.applyClient(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying tg client")
	}

	err = t.applyFolder(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying tg folder")
	}

	t.applyConfig(proj)

	return nil
}

func (t Telegram) applyClient(proj Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathsToClients, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return rerrors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	tgConnFile := patterns.TgConnFile.CopyWithNewName(
		path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), patterns.TgConnFile.Name))

	renamer.ReplaceProjectName(proj.GetName(), tgConnFile)

	proj.GetFolder().Add(
		tgConnFile,
	)

	return nil
}

func (t Telegram) applyFolder(proj Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathToServers, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	tgServer := patterns.TgServFile.Copy()

	renamer.ReplaceProjectName(proj.GetName(), tgServer)

	tgHandlerExample := patterns.TgHandlerExampleFile.Copy()

	renamer.ReplaceProjectName(proj.GetName(), tgHandlerExample)

	proj.GetFolder().Add(
		&folder.Folder{
			Name: path.Join(t.Cfg.Env.PathToServers[0], t.GetFolderName()),
			Inner: []*folder.Folder{
				tgServer,
				{
					Name:  path.Join(patterns.HandlersFolderName, patterns.VersionFolderName),
					Inner: []*folder.Folder{tgHandlerExample},
				},
			},
		},
	)

	return nil
}

func (t Telegram) applyConfig(proj Project) {
	for _, srv := range proj.GetConfig().DataSources {
		if srv.GetName() == t.GetFolderName() {
			return
		}
	}

	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, &resources.Telegram{
		Name: resources.Name(t.GetFolderName()),
	})
}
