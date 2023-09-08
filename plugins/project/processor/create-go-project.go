package processor

import (
	"os"
	"path"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/validators"
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
			go_actions.PrepareProjectStructureAction{}, // basic go project structure
			//go_actions.PrepareExamplesFoldersAction{},    // sets up examples folder for http
			go_actions.PrepareEnvironmentFoldersAction{}, // prepares environment files
			go_actions.PrepareGoConfigFolderAction{},     // config driver

			go_actions.BuildProjectAction{}, // build project in file system

			go_actions.InitGoModAction{},    // executes go mod
			go_actions.TidyAction{},         // adds/clears project initialization(api, resources) and replaces project name template with actual project name
			go_actions.FixupProjectAction{}, // fetches dependencies and formats go code

			actions.InitGit{}, // initializing and committing project as git repo
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

	proj.Cfg = config.NewEmptyConfig()

	return proj, nil
}
