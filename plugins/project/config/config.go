package config

import (
	"go.vervstack.ru/matreshka/pkg/matreshka"
)

type Config struct {
	matreshka.AppConfig

	ConfigDir  string `yaml:"-"`
	ImportPath string `yaml:"-"`
}
