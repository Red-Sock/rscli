package init

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/internal/stdio/loader"
	"github.com/Red-Sock/rscli/pkg/colors"
	"github.com/Red-Sock/rscli/plugins/project/processor"
)

var (
	projectConstructorImp = projectConstructor{
		cfg: config.GetConfig(),
		io:  stdio.StdIO{},
	}

	projectCmd = &cobra.Command{
		Use:   "project",
		Short: "Initializes project",
		Long:  `Can be used to init a project via configuration file, constructor or global config`,

		RunE: projectConstructorImp.initProject,
	}
)

func init() {
	projectCmd.Flags().StringP(nameFlag, nameFlag[:1], "", `pass a name of project with or without git pass like "rscli" or github.com/RedSock/rscli`)
	projectCmd.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)
}

type projectConstructor struct {
	cfg *config.RsCliConfig
	io  stdio.IO
}

func (p *projectConstructor) initProject(cmd *cobra.Command, _ []string) error {
	args := processor.CreateArgs{}

	// step 1: obtain name
	var err error
	args.Name, err = p.obtainName(cmd)
	if err != nil {
		return errors.Wrap(err, "error obtaining name")
	}

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(`Wonderful!!! "%s" it is!`, args.Name))

	// step 2: obtain path to project folder
	args.ProjectPath, err = p.obtainFolderPath(cmd, args.Name)
	if err != nil {
		return errors.Wrap(err, "error obtaining folder path")
	}

	err = p.buildProject(args)
	if err != nil {
		return errors.Wrap(err, "error building project")
	}
	// TODO
	return nil
}

func (p *projectConstructor) obtainName(cmd *cobra.Command) (name string, err error) {
	name = cmd.Flag(nameFlag).Value.String()

	if name == "" {
		p.io.Print(fmt.Sprintf(`
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "%s" and final result is "%s/rscli"
>`, p.cfg.DefaultProjectGitPath, p.cfg.DefaultProjectGitPath))

		name, err = p.io.GetInput()
		if err != nil {
			return "", errors.Wrap(err, "error obtaining project name")
		}
	}
	if name == "" {
		return "", errors.New("empty name")
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

	if !containsHost {
		name = p.cfg.DefaultProjectGitPath + "/" + name
	}

	err = p.validateName(name)
	if err != nil {
		return "", errors.Wrap(err, "error validating name")
	}

	return name, nil
}
func (p *projectConstructor) validateName(name string) error {
	if name == "" {
		return errors.New("no name entered")
	}

	// starting and ending ascii symbols ranges that are applicable to project name
	availableRanges := [][]int32{
		{45, 47},
		{48, 57},
		{65, 90},
		{97, 122},
	}
	for _, s := range name {
		var hasHitRange = false
		for _, r := range availableRanges {
			if s >= r[0] && s <= r[1] {
				hasHitRange = true
				break
			}
		}
		if !hasHitRange {
			return errors.New("name contains \"" + string(s) + "\" symbol")
		}
	}

	return nil
}
func (p *projectConstructor) buildProject(args processor.CreateArgs) (err error) {
	var proj *processor.Project
	proj, err = processor.CreateGoProject(args)
	if err != nil {
		return errors.Wrap(err, "error during project creation")
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

	idx := 0
	for {
		select {
		case info, ok := <-infoC:
			if !ok {
				loaders[namesToIdx[info]].Done(loader.DoneFailed)
				return nil
			}

			idx++
			loaders[namesToIdx[info]].Done(loader.DoneSuccessful)

		case buildError, ok := <-errC:
			if !ok {
				return
			}
			loaders[idx].Done(loader.DoneFailed)
			idx++
			for idx < len(loaders) {
				loaders[idx].Done(loader.DoneNotAccessed)
				idx++
			}

			return errors.Wrap(buildError, "failed to build project: ")
		}
	}
}

func (p *projectConstructor) obtainFolderPath(cmd *cobra.Command, name string) (dirPath string, err error) {
	dirPath = cmd.Flag(nameFlag).Value.String()
	if dirPath != "" {
		return dirPath, nil
	}
	dirPath, err = os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "error getting working dir")
	}

	return path.Join(dirPath, name), nil
}
