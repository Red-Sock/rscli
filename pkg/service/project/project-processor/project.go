package processor

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/actions"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/config"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/interfaces"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/validators"
	"github.com/pkg/errors"
	"os"
	"path"
)

type Action func(p interfaces.Project) error

type Validator func(p interfaces.Project) error

type Project struct {
	Name        string
	ProjectPath string
	Cfg         interfaces.Config

	Actions []Action

	validators []Validator

	F folder.Folder
}

type CreateArgs struct {
	Name        string
	CfgPath     string
	ProjectPath string
	Validators  []Validator
	Actions     []Action
}

func New(args CreateArgs) (*Project, error) {
	proj := &Project{
		Name: args.Name,
		Actions: append([]Action{
			actions.PrepareProjectStructure,   // basic project structure
			actions.PrepareConfigFolders,      // data sources and other things taken from config
			actions.PrepareAPIFolders,         // prepare servers
			actions.PrepareExamplesFolders,    // sets up examples
			actions.PrepareEnvironmentFolders, // prepares environment files

			actions.BuildConfigGoFolder, // config driver
			actions.BuildProject,        // build project in file system

			actions.InitGoMod,    // executes go mod
			actions.MoveCfg,      // moves external used config into project
			actions.FixupProject, // fetches dependencies and formats go code
		}, args.Actions...),
		F: folder.Folder{
			Name: args.Name,
		},
		validators: append(args.Validators, validators.ValidateName),
	}
	var err error
	if args.CfgPath != "" {
		proj.Cfg, err = config.NewProjectConfig(args.CfgPath)
		if err != nil {
			return proj, err
		}
	} else {
		proj.Cfg = config.NewEmptyProjectConfig()
	}

	if args.ProjectPath == "" {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return proj, errors.Wrapf(err, "error obtaining working dir")
		}
		args.ProjectPath = path.Join(wd, args.Name)
	}

	proj.ProjectPath = args.ProjectPath

	return proj, nil
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetFolder() *folder.Folder {
	return &p.F
}

func (p *Project) GetConfig() interfaces.Config {
	return p.Cfg
}

func (p *Project) GetProjectPath() string {
	return p.ProjectPath
}

func (p *Project) Build() error {
	for _, a := range p.Actions {
		if err := a(p); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) Validate() error {
	var errs []error
	for _, v := range p.validators {
		if err := v(p); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	globalErr := errors.New("error while validating the project")
	for _, e := range errs {
		globalErr = errors.Wrap(globalErr, e.Error())
	}

	return globalErr
}
