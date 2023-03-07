package structs

import "time"

type Config struct {
	AppInfo     AppInfo                `yaml:"app_info"`
	Server      map[string]interface{} `yaml:"server,omitempty"`
	DataSources map[string]interface{} `yaml:"data_sources,omitempty"`
}

type AppInfo struct {
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	StartupDuration time.Duration `yaml:"startupDuration"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]interface{}{},
		DataSources: map[string]interface{}{},
	}
}
