package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

type Redis struct {
	dependencyBase
}

func redisClient(dep dependencyBase) Dependency {
	return &Redis{
		dep,
	}
}

func (p Redis) GetFolderName() string {
	if p.Name != "" {
		return p.Name
	}

	return "redis"
}

func (p Redis) AppendToProject(proj Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return errors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Redis) applyClientFolder(proj Project) error {
	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	redisConn := patterns.RedisConnFile.CopyWithNewName(
		path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.RedisConnFile.Name))

	renamer.ReplaceProjectName(proj.GetName(), redisConn)

	proj.GetFolder().Add(redisConn)

	return nil
}

func (p Redis) applyConfig(proj Project) {
	for _, item := range proj.GetConfig().DataSources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}

	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources,
		&resources.Redis{
			Name: resources.Name(p.GetFolderName()),
			Host: "localhost",
			Port: 6379,
		})
}
