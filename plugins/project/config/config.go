package config

import (
	"github.com/godverv/matreshka"
)

type Config struct {
	*matreshka.AppConfig

	Path string `yaml:"-"`
}
