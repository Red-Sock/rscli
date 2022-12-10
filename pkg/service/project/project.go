package project

import (
	"github.com/Red-Sock/rscli/internal/utils"
	"github.com/pkg/errors"
)

func Command() []string {
	return []string{"p", "project"}
}

var (
	ErrNoArgumentsSpecifiedForFlag = errors.New("flag specified but no name was given")
	ErrFlagHasTooManyArguments     = errors.New("too many arguments specified for flag")
)

const (
	FlagAppName      = "name"
	FlagAppNameShort = "n"

	FlagCfgPath      = "cfg"
	FlagCfgPathShort = "c"
)

type Action func(p *Project) error

type Validator func(p *Project) error

type Project struct {
	Name string
	Cfg  *Config

	Actions []Action

	validators []Validator

	f Folder
}

type CreateArgs struct {
	Name       string
	CfgPath    string
	Validators []Validator
}

func NewProject(args CreateArgs) *Project {
	proj := &Project{
		Name: args.Name,
		Actions: []Action{
			prepareProjectStructure,   // basic project structure
			prepareConfigFolders,      // data sources and other things taken from config
			prepareAPIFolders,         // prepare servers
			prepareExamplesFolders,    // sets up examples
			prepareEnvironmentFolders, // prepares environment files
			buildConfigGoFolder,       // config driver
			buildProject,              // build project in file system
			initGoMod,                 // executes go mod
			moveCfg,                   // moves external used config into project
			fixupProject,
		},
		f: Folder{
			name: args.Name,
		},
		validators: append(args.Validators, ValidateName),
	}
	if args.CfgPath != "" {
		proj.Cfg = NewProjectConfig(args.CfgPath)
	}

	return proj
}

func NewProjectWithRowArgs(args []string) (*Project, error) {
	progArgs := CreateArgs{}

	flags, err := utils.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	// Define project name
	progArgs.Name, err = extractOneValueFromFlags(flags, FlagAppName, FlagAppNameShort)
	if err != nil {
		return nil, err
	}

	// Define path to configuration file
	progArgs.CfgPath, err = extractOneValueFromFlags(flags, FlagCfgPath, FlagCfgPathShort)
	if err != nil {
		return nil, err
	}
	if progArgs.CfgPath == "" {
		progArgs.CfgPath, err = findConfigPath()
		if err != nil {
			return nil, err
		}
	}

	return NewProject(progArgs), nil
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
