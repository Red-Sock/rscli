package project

import (
	"fmt"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/validators"
)

var (
	emptyNameErr = errors.New("no name entered")
)

type projectInit struct {
	io     io.IO
	config *config.RsCliConfig

	proj project.IProject
	path string
}

func newInitCmd(pi projectInit) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initializes project",
		Long:  `Can be used to init a project via configuration file, constructor or global config`,
		RunE:  pi.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return c
}

func (p *projectInit) run(_ *cobra.Command, argsIn []string) error {
	projArgs := go_project.CreateArgs{
		CfgPath: p.config.Env.PathToConfig,
	}

	// step 1: obtain name
	var err error
	projArgs.Name, err = p.obtainNameFromUser(argsIn)
	if err != nil {
		return errors.Wrap(err, "error obtaining name")
	}

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(`Wonderful!!! "%s" it is!`, projArgs.Name))

	// step 2: obtain path to project folder
	projArgs.ProjectPath = p.obtainFolderPathFromUser(projArgs.Name, argsIn)

	p.proj, err = p.buildProject(projArgs)
	if err != nil {
		return errors.Wrap(err, "error building project")
	}

	p.io.PrintlnColored(colors.ColorGreen, "Done.\nInitialized new project "+p.proj.GetName()+"\nat "+p.proj.GetProjectPath())

	return nil
}

func (p *projectInit) obtainNameFromUser(args []string) (name string, err error) {
	if len(args) > 0 {
		name = args[0]
	}

	if name == "" {
		p.io.Print(fmt.Sprintf(`
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "%s" and final result is "%s/rscli"
>`, p.config.DefaultProjectGitPath, p.config.DefaultProjectGitPath))

		name, err = p.io.GetInput()
		if err != nil {
			return "", errors.Wrap(err, "error obtaining project name")
		}
	}
	if name == "" {
		return "", emptyNameErr
	}

	if strings.HasPrefix(name, "http") {
		name = name[strings.Index(name, "://")+3:]
	}

	var containsHost bool

	// if first part (before first "/" symbol) contains dot "." -> consider it's overriding default repository
	if leftPartIdx := strings.Index(name, "/"); leftPartIdx != -1 {
		hostSeparatorIdx := strings.Index(name[:leftPartIdx], ".")
		containsHost = hostSeparatorIdx != -1
	}

	err = validators.ValidateProjectNameStr(name)
	if err != nil {
		return "", errors.Wrap(err, "error validating project name")
	}

	if !containsHost {
		name = p.config.DefaultProjectGitPath + "/" + name
	}

	name = path.Join(path.Dir(name), strings.ToLower(path.Base(name)))

	return name, nil
}

func (p *projectInit) obtainFolderPathFromUser(name string, args []string) (dirPath string) {
	if len(args) > 1 {
		return args[1]
	}

	return path.Join(p.path, path.Base(name))
}

func (p *projectInit) buildProject(args go_project.CreateArgs) (proj project.IProject, err error) {
	proj, err = go_project.CreateProject(args)
	if err != nil {
		return nil, errors.Wrap(err, "error during project creation")
	}

	p.io.Println("Starting project constructor")

	initActions := actions.InitProject(project.TypeGo)
	for _, act := range initActions {
		err = act.Do(proj)
		if err != nil {
			return nil, errors.Wrap(err, "error performing init actions")
		}
	}

	p.io.Println("Project actions performed")

	return proj, nil
}
