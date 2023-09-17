package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
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
	ok, err := containsDependency(p.Cfg, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		// TODO: RSI-141
		// p.Io.Println("already contains redis dependency")
		return nil
	}

	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoClientFolderInConfig
	}

	proj.GetFolder().Add(
		&folder.Folder{
			Name:    path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.ConnFile),
			Content: patterns.RedisConnFile,
		},
	)

	err = proj.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building pg connection folder")
	}

	return nil
}
