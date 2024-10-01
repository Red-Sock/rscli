package config

import (
	"flag"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
)

var ErrAlreadyLoaded = errors.New("config already loaded")

type Config struct {
	AppInfo matreshka.AppInfo

	DataSources DataSourcesConfig
}

var defaultConfig Config

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

func Load() (Config, error) {
	if defaultConfig.AppInfo.Name != "" {
		return defaultConfig, ErrAlreadyLoaded
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

	rootConfig, err := matreshka.ReadConfigs(cfgPath)
	if err != nil {
		return defaultConfig, errors.Wrap(err, "error reading matreshka config")
	}

	defaultConfig.AppInfo = rootConfig.AppInfo

	err = rootConfig.DataSources.ParseToStruct(defaultConfig.DataSources)
	if err != nil {
		return defaultConfig, errors.Wrap(err, "error parsing data sources to struct")
	}

	return defaultConfig, nil
}
