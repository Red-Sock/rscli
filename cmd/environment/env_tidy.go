package environment

import (
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
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

func newTidyEnvCmd() *cobra.Command {
	constr := newEnvConstructor()
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds new dependencies to existing environment. Clears unused dependencies",

		PreRunE: constr.preRun,
		RunE:    constr.runTidy,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `Path to folder with projects`)

	return c
}

func (c *envConstructor) runTidy(_ *cobra.Command, _ []string) error {
	p := ports.NewPortManager()

	errC := make(chan error)
	for _, item := range c.envProjDirs {
		go func(projName string) {
			errC <- c.tidyEnvForProject(projName, p)
		}(item.Name())
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

func (c *envConstructor) tidyEnvForProject(projName string, pm *ports.PortManager) error {
	proj, err := project.LoadProjectEnvironment(c.cfg, path.Join(c.envDirPath, projName))
	if err != nil {
		return errors.Wrap(err, "error loading environment for project "+projName)
	}

	dependencies, err := c.composePatterns.GetServiceDependencies(proj.Config)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+projName)
	}

	for _, resource := range dependencies {
		composeEnvs := resource.GetEnvs().Content()

		for _, envRow := range composeEnvs {
			if strings.HasSuffix(envRow.Name, patterns.PortSuffix) {
				var port uint64
				port, err = strconv.ParseUint(envRow.Value, 10, 16)
				if err != nil {
					return errors.Wrap(err, "error parsing .env file: port value for "+envRow.Name+" must be int but it is "+envRow.Value)
				}

				envRow.Value = strconv.FormatUint(uint64(pm.GetNextPort(uint16(port), projName)), 10)
			}

			proj.Environment.Append(envRow.Name, envRow.Value)
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
		so.Port = pm.GetNextPort(opts[idx].Port, projName)
		opts[idx] = so

		proj.Compose.Services[projName].Ports = append(
			proj.Compose.Services[projName].Ports,
			compose.AddEnvironmentBrackets(portName)+":"+strconv.Itoa(int(opts[idx].Port)))
		proj.Environment.Append(portName, strconv.Itoa(int(opts[idx].Port)))
	}

	pathToProjectEnvFile := path.Join(c.envDirPath, projName, patterns.EnvFile.Name)
	err = stdio.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectName(proj.Environment.MarshalEnv(), projName))
	if err != nil {
		return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
	}

	composeFile, err := proj.Compose.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	pathToDockerComposeFile := path.Join(c.envDirPath, projName, patterns.DockerComposeFile.Name)
	err = stdio.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectName(composeFile, projName))
	if err != nil {
		return errors.Wrap(err, "error writing docker compose file file")
	}

	return nil
}
