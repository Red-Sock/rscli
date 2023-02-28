package project

import (
	"path"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/flag"
	projsctripts "github.com/Red-Sock/rscli/plugins/project/processor"
)

var ErrProjectNoNameOrConfig = errors.New("I can't create project without name\nAlso, I like when project has git prefix.\nFor example github.com/Red-Sock/rscli whould be a great name,\nbut it is occupied :(\nTry --name or -n next time \nand use \"rscli create proj help\" for more info\n or provide a config via -c" + path.Join("path", "to", "config"))

const (
	projectHelpMsg = `
	--name, 	-n - specify the name of the project (obligatory)
	--path, 	-p - add the path to directory where project will be created (working dir by default) 
	--config, 	-c - specify the path to configuration file (otherwise empty will be created) 

`
)

const help = "help"

const (
	projectFlagName      = "name"
	projectFlagNameShort = "n"

	projectFlagPath      = "path"
	projectFlagPathShort = "p"

	projectFlagCfgPath      = "config"
	projectFlagCfgPathShort = "c"
)

func createProject(args []string) (err error) {
	argsM := flag.ParseArgs(args)

	_, ok := argsM[help]
	if ok {
		println(projectHelpMsg)
		return nil
	}
	var pArgs projsctripts.CreateArgs

	{
		pArgs.Name, err = flag.ExtractOneValueFromFlags(argsM, projectFlagName, projectFlagNameShort)
		if err != nil {
			return err
		}
	}

	{
		pArgs.ProjectPath, err = flag.ExtractOneValueFromFlags(argsM, projectFlagPath, projectFlagPathShort)
		if err != nil {
			return err
		}
	}

	{
		pArgs.CfgPath, err = flag.ExtractOneValueFromFlags(argsM, projectFlagCfgPath, projectFlagCfgPathShort)
		if err != nil {
			return err
		}
	}

	if pArgs.CfgPath == "" && pArgs.Name == "" {
		return ErrProjectNoNameOrConfig
	}

	p, err := projsctripts.CreateProject(pArgs)
	if err != nil {
		return errors.Wrapf(err, "error creating project")
	}

	err = p.Build()
	if err != nil {
		return errors.Wrapf(err, "error building project")
	}

	return nil
}
