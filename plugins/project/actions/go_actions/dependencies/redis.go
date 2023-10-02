package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type Redis struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (Redis) GetFolderName() string {
	return "redis"
}

func (p Redis) Do(proj interfaces.Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Redis) applyClientFolder(proj interfaces.Project) error {
	ok, err := containsDependency(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		return nil
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.ConnFileName),
			Content: patterns.RedisConnFile,
		},
	)

	return nil
}

func (p Redis) applyConfig(proj interfaces.Project) {
	ds := proj.GetConfig().DataSources
	if _, ok := ds[p.GetFolderName()]; ok {
		return
	}

	ds[p.GetFolderName()] = resources.Redis{
		ResourceName: "redis",
		Host:         "localhost",
		Port:         6379,
	}
}
