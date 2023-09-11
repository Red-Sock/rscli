package project

import (
	"os"
	"path"
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	pconfig "github.com/Red-Sock/rscli/plugins/project/config"
)

type Env struct {
	envDirPath  string
	Compose     *compose.Compose
	Environment *env.Container
	Config      *pconfig.Config
}

func LoadProjectEnvironment(cfg *config.RsCliConfig, pathToProjectEnv string) (p *Env, err error) {
	p = &Env{
		envDirPath: pathToProjectEnv,
	}

	err = p.fetchComposeFile()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching compose file")
	}

	err = p.fetchEnvFile()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching .env file")
	}

	err = p.fetchConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching config")
	}

	return p, nil
}

func (e *Env) Tidy(pm *ports.PortManager, composePatterns compose.PatternManager) error {
	projName := path.Base(e.envDirPath)

	err := e.tidyResources(projName, composePatterns, pm)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on resources")
	}

	err = e.tidyServerAPIs(projName, pm)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on server api")
	}

	{
		pathToProjectEnvFile := path.Join(e.envDirPath, patterns.EnvFile.Name)

		err = io.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectName(e.Environment.MarshalEnv(), projName))
		if err != nil {
			return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
		}
	}

	{
		composeFile, err := e.Compose.Marshal()
		if err != nil {
			return errors.Wrap(err, "error marshalling composer file")
		}

		pathToDockerComposeFile := path.Join(e.envDirPath, patterns.DockerComposeFile.Name)
		err = io.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectName(composeFile, projName))
		if err != nil {
			return errors.Wrap(err, "error writing docker compose file file")
		}
	}

	return nil
}

func (e *Env) tidyResources(projName string, composePatterns compose.PatternManager, pm *ports.PortManager) error {
	dataResources, err := e.Config.GetDataSourceOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining data source options")
	}

	dependencies, err := composePatterns.GetServiceDependencies(dataResources)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+e.Config.AppInfo.Name)
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

			if e.Environment.ContainsByName(newEnvName) {
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
			e.Environment.Append(newEnvName, patternEnv[idx].Value)
		}

		e.Compose.AppendService(resource.Name, resource.GetCompose())
	}

	return nil
}

func (e *Env) tidyServerAPIs(projName string, pm *ports.PortManager) error {
	opts, err := e.Config.GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining server options")
	}

	for optName := range opts {
		portName := strings.ToUpper(projName) + "_" + strings.ToUpper(opts[optName].GetName()) + "_" + patterns.PortSuffix
		portName = strings.ReplaceAll(portName,
			"__", "_")

		externalPort := opts[optName].GetPort()

		if externalPort == 0 {
			continue
		}

		e.Compose.Services[projName].Ports = append(e.Compose.Services[projName].Ports, compose.AddEnvironmentBrackets(portName)+":"+strconv.FormatUint(uint64(opts[optName].GetPort()), 10))
		e.Environment.Append(portName, strconv.FormatUint(uint64(pm.GetNextPort(opts[optName].GetPort(), portName)), 10))
	}

	return nil
}

func (e *Env) fetchComposeFile() error {
	projectEnvComposeFilePath := path.Join(e.envDirPath, patterns.DockerComposeFile.Name)
	composeFile, err := os.ReadFile(projectEnvComposeFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project env docker-compose file "+projectEnvComposeFilePath)
		}
	}

	if len(composeFile) == 0 {
		globalEnvComposeFilePath := path.Join(path.Dir(e.envDirPath), patterns.DockerComposeFile.Name)
		composeFile, err = os.ReadFile(globalEnvComposeFilePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global docker-compose file "+globalEnvComposeFilePath)
			}
		}
	}

	if len(composeFile) == 0 {
		projName := path.Base(e.envDirPath)
		composeFile = renamer.ReplaceProjectName(patterns.DockerComposeFile.Content, projName)
	}

	e.Compose, err = compose.NewComposeAssembler(composeFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}

func (e *Env) fetchEnvFile() error {
	dotEnvFilePath := path.Join(e.envDirPath, patterns.EnvFile.Name)
	envFile, err := os.ReadFile(dotEnvFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project .env file "+dotEnvFilePath)
		}
	}

	if len(envFile) == 0 {
		globalDotEnvPath := path.Join(path.Dir(e.envDirPath), patterns.EnvFile.Name)
		envFile, err = os.ReadFile(globalDotEnvPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global .env file "+globalDotEnvPath)
			}
		}
	}

	if len(envFile) == 0 {
		projName := path.Base(e.envDirPath)
		envFile = renamer.ReplaceProjectName(patterns.EnvFile.Content, projName)
	}

	e.Environment, err = env.NewEnvContainer(envFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}

// fetchConfig - searches for config in two places
// 1. in environment folder for project at ./environment/PROJ_NAME
// 2. dev.yaml file in src project (at PATH_TO_CONFIG/dev.yaml)
// if config was found by 2nd variant - it will be moved to ./environment/proj_name/dev.yaml
// and symlink will be created to it at src_proj/PATH_TO_CONFIG/dev.yaml
func (e *Env) fetchConfig(cfg *config.RsCliConfig) (err error) {
	projEnvConfigPath := path.Join(e.envDirPath, path.Base(cfg.Env.PathToConfigFolder))

	f, err := os.ReadFile(projEnvConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error ")
		}
	}

	if len(f) == 0 {
		srcProjectsDirPth := path.Dir(path.Dir(e.envDirPath))
		projName := path.Base(e.envDirPath)
		srcProjectConfigPath := path.Join(srcProjectsDirPth, projName, cfg.Env.PathToConfigFolder)

		f, err = os.ReadFile(srcProjectConfigPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "project at "+srcProjectConfigPath+" doesn't contain config")
			}
			return errors.Wrap(err, "error reading project config file")
		}

		err = os.WriteFile(projEnvConfigPath, f, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "error moving project config file to env")
		}

		err = os.RemoveAll(srcProjectConfigPath)
		if err != nil {
			return errors.Wrap(err, "error deleting config at "+srcProjectConfigPath)
		}

		err = os.Symlink(projEnvConfigPath, srcProjectConfigPath)
		if err != nil {
			return errors.Wrap(err, "error creating symlink from "+projEnvConfigPath+" to "+srcProjectConfigPath)
		}
	}

	e.Config, err = pconfig.NewConfig(f)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}

	return nil
}
