package dependencies

import (
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config/server"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type Telegram struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (t Telegram) GetFolderName() string {
	return "telegram"
}

func (t Telegram) Do(proj interfaces.Project) error {
	ok, err := containsDependency(t.Cfg, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), patterns.ConnFileName),
			Content: patterns.TgConnFile,
		},
	)

	if len(t.Cfg.Env.PathToServers) == 0 {
		return ErrNoServerFolderInConfig
	}

	for _, pth := range t.Cfg.Env.PathToServers {
		srvFold := proj.GetFolder().GetByPath(strings.Split(pth, string(os.PathSeparator))...)
		if srvFold == nil {
			continue
		}

		for _, innerFolder := range srvFold.Inner {
			if innerFolder.Name == t.GetFolderName() {
				return nil
			}
		}
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name: path.Join(t.Cfg.Env.PathToServers[0], t.GetFolderName()),
			Inner: []*folder.Folder{
				{
					Name:    patterns.TelegramServFileName,
					Content: patterns.TgServFile,
				},
				{
					Name:    path.Join(patterns.HandlersFolderName, patterns.VersionFolderName, patterns.TgHandlerFileName),
					Content: patterns.TgHandlerExampleFile,
				},
			},
		},
	)

	err = proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	ds := proj.GetConfig().Server
	if ds == nil {
		ds = make(map[string]interface{})
	}
	if _, ok = ds[t.GetFolderName()]; !ok {
		ds[t.GetFolderName()] = server.Telegram{}
	}

	proj.GetConfig().Server = ds

	return nil
}
