package environment

import (
	"context"
	stderrs "errors"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/loader"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

func newTidyEnvCmd() *cobra.Command {
	constr := newEnvConstructor()
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		RunE: constr.runTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `Path to folder with projects`)

	return c
}

func (c *envConstructor) runTidy(cmd *cobra.Command, arg []string) error {
	c.io.Println("Running rscli env tidy")

	err := c.initProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error during init of additional projects env dirs ")
	}

	err = c.fetchConstructor(cmd, arg)
	if err != nil {
		return errors.Wrap(err, "error fetching updated dirs")
	}

	portManager := ports.NewPortManager()

	progresses := make([]loader.Progress, len(c.envProjDirs))
	projectsEnvs := make([]*project.Env, len(c.envProjDirs))

	// TODO
	conflicts := make(map[uint16][]string)

	for idx := range c.envProjDirs {
		progresses[idx] = loader.NewInfiniteLoader(c.envProjDirs[idx].Name(), loader.RectSpinner())

		projName := c.envProjDirs[idx].Name()

		var proj *project.Env
		proj, err = project.LoadProjectEnvironment(c.cfg, path.Join(c.envDirPath, projName))
		if err != nil {
			return errors.Wrap(err, "error loading environment for project "+projName)
		}

		envPorts, err := proj.Environment.GetPortValues()
		if err != nil {
			return errors.Wrap(err, "error fetching ports for environment of "+projName)
		}

		for _, item := range envPorts {
			if conflictName := portManager.SaveIfNotExist(item.Value, item.Name); conflictName != "" {
				conflicts[item.Value] = []string{conflictName, item.Name}
			}
		}

		projectsEnvs[idx] = proj
	}

	done := loader.RunMultiLoader(context.Background(), c.io, progresses)
	defer func() {
		<-done()
		c.io.Println("rscli env tidy done")
	}()

	errC := make(chan error)
	for idx := range projectsEnvs {
		go func(i int) {
			err := c.tidyEnvForProject(projectsEnvs[i], portManager)
			if err != nil {
				progresses[i].Done(loader.DoneFailed)
			} else {
				progresses[i].Done(loader.DoneSuccessful)
			}

			errC <- err
		}(idx)
	}

	var errs []error
	for i := 0; i < len(c.envProjDirs); i++ {
		err, ok := <-errC
		if !ok {
			break
		}

		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}

	return stderrs.Join(errs...)
}

func (c *envConstructor) tidyEnvForProject(proj *project.Env, pm *ports.PortManager) error {
	projName := path.Base(proj.Config.AppInfo.Name)

	dataResources, err := proj.Config.GetDataSourceOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining data source options")
	}

	dependencies, err := c.composePatterns.GetServiceDependencies(dataResources)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+proj.Config.AppInfo.Name)
	}

	for _, resource := range dependencies {

		patternEnv := resource.GetEnvs().Content()

		for idx := range patternEnv {
			oldName := patternEnv[idx].Name

			newEnvName := strings.ReplaceAll(patternEnv[idx].Name,
				patterns.ResourceNameCapsPattern, strings.ToUpper(resource.Name))

			newEnvName = strings.ReplaceAll(newEnvName,
				"__", "_")

			newEnvName = string(renamer.ReplaceProjectNameStr(newEnvName, projName))

			if proj.Environment.ContainsByName(newEnvName) {
				continue
			}

			if strings.HasSuffix(newEnvName, patterns.PortSuffix) {
				var port uint64
				port, err = strconv.ParseUint(patternEnv[idx].Value, 10, 16)
				if err != nil {
					return errors.Wrap(err, "error parsing .env file: port value for "+
						newEnvName+" must be uint but it is "+
						patternEnv[idx].Value)
				}

				patternEnv[idx].Value = strconv.FormatUint(uint64(pm.GetNextPort(uint16(port), newEnvName)), 10)
			}
			resource.RenameVariable(oldName, newEnvName)
			proj.Environment.Append(newEnvName, patternEnv[idx].Value)
		}

		proj.Compose.AppendService(resource.Name, resource.GetCompose())
	}

	opts := proj.Config.GetServerOptions()

	for idx := range opts {
		if opts[idx].Port == 0 {
			continue
		}

		portName := patterns.ProjNameCapsPattern + "_" + strings.ToUpper(opts[idx].Name) + "_" + patterns.PortSuffix
		// TODO приобщить

		so := opts[idx]
		// TODO
		so.Port = pm.GetNextPort(opts[idx].Port, "projName")
		opts[idx] = so

		proj.Compose.Services[projName].Ports = append(
			proj.Compose.Services[projName].Ports,
			compose.AddEnvironmentBrackets(portName)+":"+strconv.Itoa(int(opts[idx].Port)))
		proj.Environment.Append(portName, strconv.Itoa(int(opts[idx].Port)))
	}

	pathToProjectEnvFile := path.Join(c.envDirPath, projName, patterns.EnvFile.Name)
	err = io.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectName(proj.Environment.MarshalEnv(), projName))
	if err != nil {
		return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
	}

	composeFile, err := proj.Compose.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	pathToDockerComposeFile := path.Join(c.envDirPath, projName, patterns.DockerComposeFile.Name)
	err = io.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectName(composeFile, projName))
	if err != nil {
		return errors.Wrap(err, "error writing docker compose file file")
	}

	return nil
}
