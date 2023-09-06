package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
)

type Project struct {
	pth         string
	Compose     *compose.Compose
	Environment *env.Container
	Config      *pconfig.Config
}

func LoadProjectEnvironment(cfg *config.RsCliConfig, pathToProject string) (p *Project, err error) {
	p = &Project{
		pth: pathToProject,
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

func (p *Project) fetchComposeFile() error {
	projectEnvComposeFilePath := path.Join(p.pth, patterns.DockerComposeFile.Name)
	composeFile, err := os.ReadFile(projectEnvComposeFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project env docker-compose file "+projectEnvComposeFilePath)
		}
	}

	if len(composeFile) == 0 {
		globalEnvComposeFilePath := path.Join(path.Dir(p.pth), patterns.DockerComposeFile.Name)
		composeFile, err = os.ReadFile(globalEnvComposeFilePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global docker-compose file "+globalEnvComposeFilePath)
			}
		}
	}

	if len(composeFile) == 0 {
		projName := path.Base(p.pth)
		composeFile = renamer.ReplaceProjectName(patterns.DockerComposeFile.Content, projName)
	}

	p.Compose, err = compose.NewComposeAssembler(composeFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}

func (p *Project) fetchEnvFile() error {
	dotEnvFilePath := path.Join(p.pth, patterns.EnvFile.Name)
	envFile, err := os.ReadFile(dotEnvFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project .env file "+dotEnvFilePath)
		}
	}

	if len(envFile) == 0 {
		globalDotEnvPath := path.Join(path.Dir(p.pth), patterns.DockerComposeFile.Name)
		envFile, err = os.ReadFile(globalDotEnvPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global .env file "+globalDotEnvPath)
			}
		}
	}

	if len(envFile) == 0 {
		projName := path.Base(p.pth)
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
	projEnvConfigPath := path.Join(p.pth, path.Base(cfg.Env.PathToConfig))

	f, err := os.ReadFile(projEnvConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error ")
		}
	}

	if len(f) == 0 {
		srcProjectsDirPth := path.Dir(path.Dir(p.pth))
		projName := path.Base(p.pth)
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
