package config

import (
	"flag"
	"time"

	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
)

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

var (
	ErrNotFound    = errors.New("no such key")
	ErrCannotParse = errors.New("couldn't parse value")
)

type Config interface {
	GetInt(key string) (out int)
	GetString(key string) (out string)
	GetBool(key string) (out bool)
	GetDuration(key string) (out time.Duration)

	TryGetInt(key string) (out int, err error)
	TryGetString(key string) (out string, err error)
	TryGetBool(key string) (out bool, err error)
	TryGetDuration(key string) (t time.Duration, err error)
}

var defaultConfig Config

func GetConfig() (Config, error) {
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

	cfg, err := matreshka.ReadConfig(cfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading matreshka config")
	}

	return cfg, nil
}

func GetInt(key string) (out int) {
	panic("not implemented")
}
func GetString(key string) (out string) {
	panic("not implemented")
}
func GetBool(key string) (out bool) {
	panic("not implemented")
}
func GetDuration(key string) (out time.Duration) {
	panic("not implemented")
}

func TryGetInt(key string) (out int, err error) {
	panic("not implemented")
}
func TryGetString(key string) (out string, err error) {
	panic("not implemented")
}
func TryGetBool(key string) (out bool, err error) {
	panic("not implemented")
}
func TryGetDuration(key string) (t time.Duration, err error) {
	panic("not implemented")
}
