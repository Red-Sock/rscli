package processor

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/validators"
)

type Action func(p interfaces.Project) error

type Validator func(p interfaces.Project) error

type Project struct {
	Name        string
	ModName     string
	ProjectPath string
	Cfg         interfaces.ProjectConfig

	Actions []Action

	F folder.Folder

	RscliVersion interfaces.Version

	validators []Validator
}

type CreateArgs struct {
	Name        string
	CfgPath     string
	ProjectPath string
	Validators  []Validator
	Actions     []Action
}

func CreateProject(args CreateArgs) (*Project, error) {
	proj := &Project{
		Name: args.Name,
		Actions: append([]Action{
			actions.PrepareProjectStructure,   // basic project structure
			actions.PrepareConfigFolders,      // data sources and other things taken from config
			actions.PrepareExamplesFolders,    // sets up examples
			actions.PrepareEnvironmentFolders, // prepares environment files

			actions.BuildConfigGoFolder, // config driver
			actions.BuildProject,        // build project in file system

			actions.InitGoMod,    // executes go mod
			actions.MoveCfg,      // moves external used config into project
			actions.Tidy,         // adds/clears project initialization(api, resources) and replaces project name template with actual project name
			actions.FixupProject, // fetches dependencies and formats go code
			actions.InitGit,      // initializing and committing project as git repo
		}, args.Actions...),
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

	if args.Name == "" {
		proj.ModName, err = proj.Cfg.ExtractName()
		if err != nil {
			return nil, err
		}
		args.Name = proj.Name
	}

	proj.Name = proj.ModName

	if projectNameStartIdx := strings.LastIndex(proj.ModName, "/"); projectNameStartIdx != -1 {
		proj.Name = proj.Name[projectNameStartIdx+1:]
	}

	if args.ProjectPath == "" {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return proj, errors.Wrapf(err, "error obtaining working dir")
		}
		args.ProjectPath = path.Join(wd, proj.Name)
	}

	proj.F = folder.Folder{
		Name: proj.Name,
	}

	proj.ProjectPath = args.ProjectPath

	return proj, nil
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetProjectModName() string {
	return p.ModName
}

func (p *Project) GetFolder() *folder.Folder {
	return &p.F
}

func (p *Project) GetConfig() interfaces.ProjectConfig {
	return p.Cfg
}

func (p *Project) GetProjectPath() string {
	return p.ProjectPath
}

func (p *Project) Build() (err error) {
	for _, a := range p.Actions {
		if err = a(p); err != nil {
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

func (p *Project) GetVersion() interfaces.Version {
	return p.RscliVersion
}

func (p *Project) SetVersion(version interfaces.Version) {
	p.RscliVersion = version
}
