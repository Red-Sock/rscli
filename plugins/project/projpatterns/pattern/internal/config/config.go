package config

import (
	"time"

	"github.com/godverv/matreshka"
)

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

type Config interface {
	AppInfo() matreshka.AppInfo
	Api() matreshka.Servers
	Resources() matreshka.Resources

	GetInt(key string) (out int)
	GetString(key string) (out string)
	GetBool(key string) (out bool)
	GetDuration(key string) (out time.Duration)

	TryGetInt(key string) (out int, err error)
	TryGetString(key string) (out string, err error)
	TryGetBool(key string) (out bool, err error)
	TryGetDuration(key string) (t time.Duration, err error)
}

var defaultConfig *config

type config struct {
	matreshka.AppConfig
}

func GetConfig() Config {
	return defaultConfig
}

func (c *config) AppInfo() matreshka.AppInfo {
	return c.AppConfig.AppInfo
}

func (c *config) Api() matreshka.Servers {
	return c.AppConfig.Servers
}

func (c *config) Resources() matreshka.Resources {
	return c.AppConfig.Resources
}
