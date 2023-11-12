package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
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
	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding dependency path")
	}

	if ok {
		return nil
	}

	redisConn := projpatterns.RedisConnFile.CopyWithNewName(
		path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), projpatterns.RedisConnFile.Name))

	go_actions.ReplaceProjectName(proj.GetName(), redisConn)

	proj.GetFolder().Add(redisConn)

	return nil
}

func (p Redis) applyConfig(proj interfaces.Project) {
	for _, item := range proj.GetConfig().Resources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}

	proj.GetConfig().Resources = append(proj.GetConfig().Resources,
		&resources.Redis{
			Name: resources.Name(p.GetFolderName()),
			Host: "localhost",
			Port: 6379,
		})
}
