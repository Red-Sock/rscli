package config

import (
	"flag"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
)

var ErrAlreadyLoaded = errors.New("config already loaded")

var defaultConfig matreshka.AppConfig

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

func Load() (matreshka.Config, error) {
	if defaultConfig.AppInfo.Name != "" {
		return &defaultConfig, ErrAlreadyLoaded
	}

	var cfgPath string
	var isDevBuild bool

	flag.StringVar(&cfgPath, "config", "", "Path to configuration file")
	flag.BoolVar(&isDevBuild, "dev", false, "Flag turns on a dev config at ./config/dev.yaml")
	flag.Parse()

	if cfgPath == "" {
		if isDevBuild {
			cfgPath = devConfigPath
		} else {
			cfgPath = prodConfigPath
		}
	}
	var err error
	defaultConfig, err = matreshka.ReadConfigs(cfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading matreshka config")
	}

	return &defaultConfig, nil
}
