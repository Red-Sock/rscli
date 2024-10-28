package project

import (
	"os"
	"path"
	"time"

	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

const (
	defaultVersion  = "v0.0.1"
	startupDuration = time.Second * 10
)

type CreateArgs struct {
	Name        string
	CfgPath     string
	ProjectPath string

	Type Type
}

func CreateProject(args CreateArgs) (*Project, error) {
	proj := &Project{
		Name: args.Name,

		ProjType: args.Type,
	}

	if args.ProjectPath == "" {
		var wd string
		wd, err := os.Getwd()
		if err != nil {
			return proj, errors.Wrapf(err, "error obtaining working dir")
		}

		args.ProjectPath = path.Join(wd, proj.Name)
	}

	proj.Path = args.ProjectPath

	if args.CfgPath == "" {
		args.CfgPath = rscliconfig.GetConfig().Env.PathToConfig
	}

	proj.Cfg = &config.Config{
		AppConfig: matreshka.AppConfig{
			AppInfo: matreshka.AppInfo{
				Name:            proj.GetName(),
				Version:         defaultVersion,
				StartupDuration: startupDuration,
			},
		},
		ConfigDir: path.Join(proj.GetProjectPath(), args.CfgPath),
	}

	return proj, nil
}
