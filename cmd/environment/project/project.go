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
	projpatterns "github.com/Red-Sock/rscli/plugins/project/patterns"
)

var ErrNoConfig = errors.New("no config found")

type envResourcePattern interface {
	GetByName(envName string) (value string)
}

type Env struct {
	envDirPath  string
	Compose     *compose.Compose
	Environment *env.Container
	Config      *pconfig.Config

	environmentResourcePatterns envResourcePattern
}

func LoadProjectEnvironment(cfg *config.RsCliConfig, envResourcePattern envResourcePattern, pathToProjectEnv string) (p *Env, err error) {
	p = &Env{
		envDirPath:                  pathToProjectEnv,
		environmentResourcePatterns: envResourcePattern,
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

	e.tidyEnvFile()

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

func (e *Env) tidyEnvFile() {
	for _, envVar := range e.Environment.Content() {
		if envVar.Name == "" || envVar.Name[0] == '#' {
			e.Environment.Remove(envVar.Name)
		}
	}
}

func (e *Env) tidyResources(projName string, composePatterns compose.PatternManager, pm *ports.PortManager) error {
	dependencies, err := composePatterns.GetServiceDependencies(e.Config)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+e.Config.AppInfo.Name)
	}

	for _, resource := range dependencies {

		patternEnv := resource.GetEnvs().Content()

		for idx := range patternEnv {

			newEnvName := strings.ReplaceAll(patternEnv[idx].Name,
				patterns.ResourceNameCapsPattern, strings.ToUpper(resource.Name))

			{
				newEnvName = strings.ReplaceAll(newEnvName,
					"__", "_")
				newEnvName = renamer.ReplaceProjectNameStr(newEnvName, projName)
			}

			if e.Environment.ContainsByName(newEnvName) {
				continue
			}

			newEnvValue := e.environmentResourcePatterns.GetByName(newEnvName)

			if strings.HasSuffix(newEnvName, patterns.PortSuffix) {
				var port uint64
				port, err = strconv.ParseUint(patternEnv[idx].Value, 10, 16)
				if err != nil {
					return errors.Wrap(err, "error parsing .env file: port value for "+
						newEnvName+" must be uint but it is "+
						patternEnv[idx].Value)
				}

				newEnvValue = strconv.FormatUint(uint64(pm.GetNextPort(uint16(port), newEnvName)), 10)
			} else {
				newEnvValue = renamer.ReplaceProjectNameStr(newEnvValue, projName)
			}

			resource.RenameVariable(patternEnv[idx].Name, newEnvName)

			e.Environment.Append(newEnvName, newEnvValue)
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

	service, ok := e.Compose.Services[strings.ReplaceAll(projName, "-", "_")]
	if !ok {
		service = e.Compose.Services[patterns.ProjNamePattern]
	}

	for optName := range opts {
		portName := strings.ToUpper(projName) + "_" + strings.ToUpper(opts[optName].GetName()) + "_" + patterns.PortSuffix
		portName = strings.ReplaceAll(portName,
			"__", "_")

		externalPort := opts[optName].GetPort()

		if externalPort == 0 {
			continue
		}

		composePort := compose.AddEnvironmentBrackets(portName) + ":" + strconv.FormatUint(uint64(opts[optName].GetPort()), 10)
		portExists := false
		for _, item := range service.Ports {
			if item == composePort {
				portExists = true
				break
			}
		}
		if !portExists {
			service.Ports = append(service.Ports, composePort)
		}

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
func (e *Env) fetchConfig(cfg *config.RsCliConfig) error {
	f, err := e.findEnvConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "error finding environment config")
	}

	e.Config, err = pconfig.NewConfig(f)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}
	return nil
}

func (e *Env) findEnvConfig(cfg *config.RsCliConfig) ([]byte, error) {
	// trying to find env.yaml file in env folder
	envConfigPath := path.Join(e.envDirPath, projpatterns.EnvConfigYamlFile)

	f, err := os.ReadFile(envConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrap(err, "error reading environment config file")
		}
	}
	if len(f) != 0 {
		return f, nil
	}

	srcProjectsDirPth := path.Dir(path.Dir(e.envDirPath))
	projName := path.Base(e.envDirPath)
	projEnvConfigPath := path.Join(srcProjectsDirPth, projName, path.Dir(cfg.Env.PathToConfig), projpatterns.EnvConfigYamlFile)

	// trying to find env.yaml file in project folder (might be left from previous "rscli env" use)
	f, err = os.ReadFile(projEnvConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrap(err, "error reading environment config file in project")
		}
	}

	if len(f) != 0 {
		err = os.Link(projEnvConfigPath, envConfigPath)
		if err != nil {
			return nil, errors.Wrap(err, "error creating hardlink from "+projEnvConfigPath+" to "+envConfigPath)
		}
		return f, nil
	}

	// trying to find default config file in project folder (might be left from previous "rscli env" use)

	srcProjectConfigPath := path.Join(srcProjectsDirPth, projName, cfg.Env.PathToConfig)

	f, err = os.ReadFile(srcProjectConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrap(err, "project at "+srcProjectConfigPath+" doesn't contain config")
		}
		return nil, errors.Wrap(err, "error reading project config file")
	}

	if len(f) == 0 {
		return nil, errors.Wrap(ErrNoConfig, "no config found")
	}

	err = os.WriteFile(envConfigPath, f, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "error creating environment config at "+envConfigPath)
	}

	err = os.Link(envConfigPath, projEnvConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error creating hardlink from "+envConfigPath+" to "+projEnvConfigPath)
	}

	return f, nil
}
