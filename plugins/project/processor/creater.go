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

type Action interface {
	Do(p interfaces.Project) error
	NameInAction() string
}

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
			actions.PrepareProjectStructureAction{},   // basic project structure
			actions.PrepareExamplesFoldersAction{},    // sets up examples
			actions.PrepareEnvironmentFoldersAction{}, // prepares environment files

			actions.BuildConfigGoFolderAction{}, // config driver
			actions.BuildProjectAction{},        // build project in file system

			actions.InitGoModAction{},    // executes go mod
			actions.MoveCfgAction{},      // moves external used config into project
			actions.TidyAction{},         // adds/clears project initialization(api, resources) and replaces project name template with actual project name
			actions.FixupProjectAction{}, // fetches dependencies and formats go code
			actions.InitGit{},            // initializing and committing project as git repo
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

func (p *Project) Build() (<-chan string, <-chan error) {
	progressCh := make(chan string, len(p.Actions))
	errCh := make(chan error)

	go func() {
		for _, a := range p.Actions {
			progressCh <- a.NameInAction()
			if err := a.Do(p); err != nil {
				close(progressCh)
				errCh <- err
				close(errCh)
				return
			}
		}
		close(progressCh)
		close(errCh)
	}()

	return progressCh, errCh
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
