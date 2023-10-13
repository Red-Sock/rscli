package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/makefile"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/config"
	pconfig "github.com/Red-Sock/rscli/plugins/project/config"
	projpatterns "github.com/Red-Sock/rscli/plugins/project/patterns"
)

// fetchConfig - searches for config in two places
// 1. in environment folder for project at ./environment/PROJ_NAME
// 2. dev.yaml file in src project (at PATH_TO_CONFIG/dev.yaml)
// if config was found by 2nd variant - it will be moved to ./environment/proj_name/dev.yaml
// and symlink will be created to it at src_proj/PATH_TO_CONFIG/dev.yaml
func (e *Env) fetchConfig(cfg *config.RsCliConfig) error {
	confPath, err := e.findEnvConfig(cfg)
	if err != nil {
		if !errors.Is(err, ErrNoConfig) {
			return errors.Wrap(err, "error finding environment config")
		}
	}
	{
		srcProjectsDirPth := path.Dir(path.Dir(e.envDirPath))
		projName := path.Base(e.envDirPath)
		projEnvConfigPath := path.Join(srcProjectsDirPth, projName, path.Dir(cfg.Env.PathToConfig), projpatterns.EnvConfigYamlFile)

		_, err = os.Stat(projEnvConfigPath)
		if err != nil {
			err = os.Symlink(confPath, projEnvConfigPath)
			if err != nil {
				return errors.Wrap(err, "error creating symlink")
			}
		}
	}
	e.Config, err = pconfig.ReadConfig(confPath)
	if err != nil {
		return errors.Wrap(err, "error parsing config")
	}

	return nil
}

func (e *Env) findEnvConfig(cfg *config.RsCliConfig) (string, error) {
	// trying to find env.yaml file in env folder
	envConfigPath := path.Join(e.envDirPath, projpatterns.EnvConfigYamlFile)

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

	srcProjectsDirPth := path.Dir(path.Dir(e.envDirPath))
	projName := path.Base(e.envDirPath)
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

func (e *Env) fetchMakeFile() (err error) {
	e.Makefile, err = makefile.ReadMakeFile(path.Join(e.envDirPath, patterns.Makefile.Name))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error getting makefile")
		}

		e.Makefile = makefile.MewEmptyMakefile()
	}

	return nil
}
