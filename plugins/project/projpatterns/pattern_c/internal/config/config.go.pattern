package config

import (
	"time"

	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/api"
	"github.com/godverv/matreshka/resources"
)

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

type Config interface {
	AppInfo() matreshka.AppInfo
	Api() API
	Resources() Resource

	GetInt(key string) (out int)
	GetString(key string) (out string)
	GetBool(key string) (out bool)
	GetDuration(key string) (out time.Duration)

	TryGetInt(key string) (out int, err error)
	TryGetString(key string) (out string, err error)
	TryGetBool(key string) (out bool, err error)
	TryGetDuration(key string) (t time.Duration, err error)

	GetMatreshka() *matreshka.AppConfig
}
type API interface {
	REST(name string) (*api.Rest, error)
	GRPC(name string) (*api.GRPC, error)
}
type Resource interface {
	Postgres(name string) (*resources.Postgres, error)
	Telegram(name string) (*resources.Telegram, error)
	Redis(name string) (*resources.Redis, error)
	GRPC(name string) (*resources.GRPC, error)
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

func (c *config) Api() API {
	return &c.AppConfig.Servers
}

func (c *config) Resources() Resource {
	return &c.AppConfig.Resources
}

func (c *config) GetMatreshka() *matreshka.AppConfig {
	return &c.AppConfig
}
