package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type envConfig struct {
	*matreshka.AppConfig
}

// fetch - searches for config in two places
// 1. in environment folder for project at ./environment/PROJ_NAME
// 2. dev.yaml file in src project (at PATH_TO_CONFIG/dev.yaml)
// if config was found by 2nd variant - it will be moved to ./environment/proj_name/dev.yaml
// and symlink will be created to it at src_proj/PATH_TO_CONFIG/dev.yaml
func (e *envConfig) fetch(cfg *config.RsCliConfig, pathToProjectEnv, pathToProject string) error {
	confPath, err := e.findEnvConfig(cfg, pathToProjectEnv)
	if err != nil {
		if !errors.Is(err, ErrNoConfig) {
			return errors.Wrap(err, "error finding environment config")
		}
	}

	{
		srcProjectsDirPth := path.Dir(path.Dir(pathToProjectEnv))
		projName := path.Base(pathToProjectEnv)
		projEnvConfigPath := path.Join(srcProjectsDirPth, projName, path.Dir(cfg.Env.PathToConfig), projpatterns.EnvConfigYamlFile)

		_, err = os.Stat(projEnvConfigPath)
		if err != nil {
			err = os.Symlink(confPath, projEnvConfigPath)
			if err != nil {
				return errors.Wrap(err, "error creating symlink")
			}
		}
	}

	e.AppConfig, err = matreshka.ReadConfig(confPath)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}

	projConfig, err := project.LoadProjectConfig(pathToProject, cfg)
	if err != nil {
		return nil
	}

	idx := 0
	for i := len(e.AppConfig.Resources) - 1; i > 0; i-- {
		var contains bool
		for _, projectConfig := range projConfig.Resources {
			if projectConfig.GetName() == e.AppConfig.Resources[i].GetName() {
				contains = true
				break
			}
		}
		if !contains {
			e.AppConfig.Resources[idx], e.AppConfig.Resources[i] = e.AppConfig.Resources[i], e.AppConfig.Resources[idx]
			idx++
		}

	}

	for _, pc := range projConfig.Resources {
		var contains bool
		for _, ec := range e.AppConfig.Resources {
			if pc.GetName() == ec.GetName() {
				contains = true
				break
			}
		}

		if !contains {
			e.AppConfig.Resources = append(e.AppConfig.Resources, pc)
		}
	}

	return nil
}

func (e *envConfig) findEnvConfig(cfg *config.RsCliConfig, pathToProjectEnv string) (string, error) {
	// trying to find env.yaml file in env folder
	envConfigPath := path.Join(pathToProjectEnv, projpatterns.EnvConfigYamlFile)

	s, err := os.Stat(envConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", errors.Wrap(err, "error reading environment config file")
		}
	} else {
		if !s.IsDir() {
			return envConfigPath, nil
		}
	}

	srcProjectsDirPth := path.Dir(path.Dir(pathToProjectEnv))
	projName := path.Base(pathToProjectEnv)
	projEnvConfigPath := path.Join(srcProjectsDirPth, projName, path.Dir(cfg.Env.PathToConfig), projpatterns.EnvConfigYamlFile)

	// trying to find env.yaml file in project folder (might be left from previous "rscli env" use)
	stat, err := os.Stat(projEnvConfigPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", errors.Wrap(err, "error reading environment config file in project")
		}
	} else {
		if !stat.IsDir() {
			return envConfigPath, nil
		}
	}

	srcProjectConfigPath := path.Join(srcProjectsDirPth, projName, cfg.Env.PathToConfig)

	f, err := os.ReadFile(srcProjectConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.Wrap(ErrNoConfig, "project at "+srcProjectConfigPath+" doesn't contain config")
		}
		return "", errors.Wrap(err, "error reading project config file")
	}

	err = os.WriteFile(envConfigPath, f, os.ModePerm)
	if err != nil {
		return "", errors.Wrap(err, "error writing env config from src project")
	}

	return envConfigPath, nil
}
