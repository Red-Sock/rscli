package config

import (
	"flag"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
)

var ErrAlreadyLoaded = errors.New("config already loaded")

func Load() (Config, error) {
	if defaultConfig != nil {
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

	cfg, err := matreshka.ReadConfigs(cfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading matreshka config")
	}

	defaultConfig = &config{AppConfig: *cfg}

	return defaultConfig, nil
}
