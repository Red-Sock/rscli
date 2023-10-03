package project

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/internal/io/loader"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/validators"
)

var (
	emptyNameErr = errors.New("no name entered")
)

type projectInit struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
	path string
}

func newInitCmd(pi projectInit) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initializes project",
		Long:  `Can be used to init a project via configuration file, constructor or global config`,

		RunE: pi.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(nameFlag, nameFlag[:1], "", `name of project with or without git pass like "rscli" or github.com/RedSock/rscli`)
	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectInit) run(cmd *cobra.Command, _ []string) error {
	args := project.CreateArgs{}

	// step 1: obtain name
	var err error
	args.Name, err = p.obtainNameFromUser(cmd)
	if err != nil {
		return errors.Wrap(err, "error obtaining name")
	}

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(`Wonderful!!! "%s" it is!`, args.Name))

	// step 2: obtain path to project folder
	args.ProjectPath = p.obtainFolderPathFromUser(cmd, args.Name)

	p.proj, err = p.buildProject(args)
	if err != nil {
		return errors.Wrap(err, "error building project")
	}

	p.io.PrintlnColored(colors.ColorGreen, "Done.\nInitialized new project "+p.proj.GetName()+"\nat "+p.proj.GetProjectPath())

	return nil
}

func (p *projectInit) obtainNameFromUser(cmd *cobra.Command) (name string, err error) {
	name = cmd.Flag(nameFlag).Value.String()

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

	return name, nil
}

func (p *projectInit) obtainFolderPathFromUser(cmd *cobra.Command, name string) (dirPath string) {
	dirPath = cmd.Flag(pathFlag).Value.String()
	if dirPath != "" {
		return dirPath
	}

	dirPath = path.Join(p.path, path.Base(name))

	return dirPath
}

func (p *projectInit) buildProject(args project.CreateArgs) (proj *project.Project, err error) {
	proj, err = project.CreateGoProject(args)
	if err != nil {
		return nil, errors.Wrap(err, "error during project creation")
	}

	actionNames := proj.GetActionNames()

	loaders := make([]loader.Progress, 0, len(actionNames))
	namesToIdx := make(map[string]int, len(loaders))

	for idx, actionName := range actionNames {
		loaders = append(loaders, loader.NewInfiniteLoader(actionName, loader.RectSpinner()))
		namesToIdx[actionName] = idx
	}

	infoC, errC := proj.Build()

	p.io.Println("Starting project constructor")

	doneF := loader.RunMultiLoader(context.TODO(), p.io, loaders)
	defer func() {
		<-doneF()
	}()

	currentProcessIdx := 0

	fail := func() {
		loaders[currentProcessIdx].Done(loader.DoneFailed)
		currentProcessIdx++
		for currentProcessIdx < len(loaders) {
			loaders[currentProcessIdx].Done(loader.DoneNotAccessed)
			currentProcessIdx++
		}
	}

	for {
		select {
		case info, ok := <-infoC:
			if !ok {
				fail()
				return proj, nil
			}

			currentProcessIdx++
			loaders[namesToIdx[info]].Done(loader.DoneSuccessful)

		case buildError, ok := <-errC:
			if !ok {
				return proj, nil
			}

			fail()
			return nil, errors.Wrap(buildError, "failed to build project: ")
		}
	}
}
