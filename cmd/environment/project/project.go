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
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
)

type Project struct {
	envDirPath  string
	Compose     *compose.Compose
	Environment *env.Container
	Config      *pconfig.Config
}

func LoadProjectEnvironment(cfg *config.RsCliConfig, pathToProjectEnv string) (p *Project, err error) {
	p = &Project{
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

func (p *Project) Tidy(pm *ports.PortManager, composePatterns compose.PatternManager) error {
	projName := path.Base(p.envDirPath)
	err := p.tidyResources(projName, composePatterns, pm)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on resources")
	}

	p.tidyServerAPIs(projName, pm)

	pathToProjectEnvFile := path.Join(p.envDirPath, patterns.EnvFile.Name)
	err = stdio.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectName(p.Environment.MarshalEnv(), projName))
	if err != nil {
		return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
	}

	composeFile, err := p.Compose.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	pathToDockerComposeFile := path.Join(p.envDirPath, patterns.DockerComposeFile.Name)
	err = stdio.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectName(composeFile, projName))
	if err != nil {
		return errors.Wrap(err, "error writing docker compose file file")
	}

	return nil
}

func (p *Project) tidyResources(projName string, composePatterns compose.PatternManager, pm *ports.PortManager) error {
	dependencies, err := composePatterns.GetServiceDependencies(p.Config)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+projName)
	}

	for _, resource := range dependencies {
		composeEnvs := resource.GetEnvs().Content()

		for _, envRow := range composeEnvs {
			if strings.HasSuffix(envRow.Name, patterns.PortSuffix) {

				if p.Environment.Contains(env.Variable{
					Name:  envRow.Name,
					Value: envRow.Value,
				}) {
					continue
				}

				var port uint64
				port, err = strconv.ParseUint(envRow.Value, 10, 16)
				if err != nil {
					return errors.Wrap(err, "error parsing .env file: port value for "+envRow.Name+" must be int but it is "+envRow.Value)
				}

				envRow.Value = strconv.FormatUint(uint64(pm.GetNextPort(uint16(port), projName)), 10)
			}

			p.Environment.Append(envRow.Name, envRow.Value)
		}

		p.Compose.AppendService(resource.Name, resource.GetCompose())
	}

	return nil
}

func (p *Project) tidyServerAPIs(projName string, pm *ports.PortManager) {
	serverAPIs := p.Config.GetServerOptions()

	for idx := range serverAPIs {
		if serverAPIs[idx].Port == 0 {
			continue
		}

		portName := patterns.ProjNameCapsPattern + "_" + strings.ToUpper(serverAPIs[idx].Name) + "_" + patterns.PortSuffix
		portValue := strconv.Itoa(int(serverAPIs[idx].Port))
		if p.Environment.Contains(env.Variable{
			Name:  portName,
			Value: portValue,
		}) {
			continue
		}

		// TODO приобщить
		so := serverAPIs[idx]
		so.Port = pm.GetNextPort(serverAPIs[idx].Port, projName)
		serverAPIs[idx] = so

		portValue = strconv.Itoa(int(serverAPIs[idx].Port))

		p.Compose.Services[projName].Ports = append(
			p.Compose.Services[projName].Ports,
			compose.AddEnvironmentBrackets(portName)+":"+portValue)

		p.Environment.Append(portName, portValue)
	}
}

func (p *Project) fetchComposeFile() error {
	projectEnvComposeFilePath := path.Join(p.envDirPath, patterns.DockerComposeFile.Name)
	composeFile, err := os.ReadFile(projectEnvComposeFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project env docker-compose file "+projectEnvComposeFilePath)
		}
	}

	if len(composeFile) == 0 {
		globalEnvComposeFilePath := path.Join(path.Dir(p.envDirPath), patterns.DockerComposeFile.Name)
		composeFile, err = os.ReadFile(globalEnvComposeFilePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global docker-compose file "+globalEnvComposeFilePath)
			}
		}
	}

	if len(composeFile) == 0 {
		projName := path.Base(p.envDirPath)
		composeFile = renamer.ReplaceProjectName(patterns.DockerComposeFile.Content, projName)
	}

	p.Compose, err = compose.NewComposeAssembler(composeFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}

func (p *Project) fetchEnvFile() error {
	dotEnvFilePath := path.Join(p.envDirPath, patterns.EnvFile.Name)
	envFile, err := os.ReadFile(dotEnvFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project .env file "+dotEnvFilePath)
		}
	}

	if len(envFile) == 0 {
		globalDotEnvPath := path.Join(path.Dir(p.envDirPath), patterns.DockerComposeFile.Name)
		envFile, err = os.ReadFile(globalDotEnvPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global .env file "+globalDotEnvPath)
			}
		}
	}

	if len(envFile) == 0 {
		projName := path.Base(p.envDirPath)
		envFile = renamer.ReplaceProjectName(patterns.EnvFile.Content, projName)
	}

	p.Environment, err = env.NewEnvContainer(envFile)
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
func (p *Project) fetchConfig(cfg *config.RsCliConfig) (err error) {
	projEnvConfigPath := path.Join(p.envDirPath, path.Base(cfg.Env.PathToConfig))

	f, err := os.ReadFile(projEnvConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error ")
		}
	}

	if len(f) == 0 {
		srcProjectsDirPth := path.Dir(path.Dir(p.envDirPath))
		projName := path.Base(p.envDirPath)
		srcProjectConfigPath := path.Join(srcProjectsDirPth, projName, cfg.Env.PathToConfig)

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

	p.Config, err = pconfig.NewConfig(f)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}

	return nil
}
