package config

import (
	"go.vervstack.ru/matreshka"
)

type Config struct {
	matreshka.AppConfig

	ConfigDir  string `yaml:"-"`
	ImportPath string `yaml:"-"`
}
