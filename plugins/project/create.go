package project

import (
	"os"
	"path"
	"time"

	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/validators"
)

const (
	defaultVersion  = "0.0.1"
	startupDuration = time.Second * 10
)

type CreateArgs struct {
	Name        string
	CfgPath     string
	ProjectPath string
	Validators  []Validator
	Actions     []actions.Action
}

func CreateGoProject(args CreateArgs) (*Project, error) {
	proj := &Project{
		Name: args.Name,
		Actions: append([]actions.Action{
			go_actions.PrepareProjectStructureAction{},   // basic go project structure
			go_actions.PrepareEnvironmentFoldersAction{}, // prepares environment files
			go_actions.PrepareGoConfigFolderAction{},

			go_actions.BuildProjectAction{}, // build project in file system

			go_actions.InitGoModAction{}, // executes go mod
			go_actions.TidyAction{},      // adds/clears project initialization(api, resources) and replaces project name template with actual project name
			go_actions.FormatAction{},    // fetches dependencies and formats go code

			git.InitGit{}, // initializing and committing project as git repo
		}, args.Actions...),
		validators: append(args.Validators, validators.ValidateProjectName),
	}

	if args.ProjectPath == "" {
		var wd string
		wd, err := os.Getwd()
		if err != nil {
			return proj, errors.Wrapf(err, "error obtaining working dir")
		}

		args.ProjectPath = path.Join(wd, proj.Name)
	}

	proj.ProjectPath = args.ProjectPath

	proj.root = folder.Folder{
		Name: proj.ProjectPath,
	}

	proj.Cfg = &config.Config{
		AppConfig: &matreshka.AppConfig{
			AppInfo: matreshka.AppInfo{
				Name:            proj.GetName(),
				Version:         defaultVersion,
				StartupDuration: startupDuration,
			},
		},
		Path: path.Join(proj.GetProjectPath(), rscliconfig.GetConfig().Env.PathToConfig),
	}

	return proj, nil
}
