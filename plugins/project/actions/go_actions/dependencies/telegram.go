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

	if ok {
		// TODO: RSI-141
		// t.Io.Println("already contains pg dependency")
		return nil
	}

	if len(t.Cfg.Env.PathsToClients) == 0 {
		return ErrNoClientFolderInConfig
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), patterns.ConnFileName),
			Content: patterns.TgConnFile,
		},
	)

	err = proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	ds := proj.GetConfig().Server
	if _, ok = ds[t.GetFolderName()]; !ok {
		ds[t.GetFolderName()] = server.Telegram{}
	}

	return nil
}
