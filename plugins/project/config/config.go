package config

import (
	"go.verv.tech/matreshka"
)

type Config struct {
	matreshka.AppConfig

	ConfigDir  string `yaml:"-"`
	ImportPath string `yaml:"-"`
}
