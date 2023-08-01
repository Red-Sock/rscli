package configstructs

import "time"

type Config struct {
	AppInfo     AppInfo                `yaml:"app_info"`
	Server      map[string]interface{} `yaml:"server,omitempty"`
	DataSources map[string]interface{} `yaml:"data_sources,omitempty"`
}

type AppInfo struct {
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	StartupDuration time.Duration `yaml:"startup_duration"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]interface{}{},
		DataSources: map[string]interface{}{},
	}
}

type ServerOptions struct {
	Name        string
	Port        uint16 `yaml:"port"`
	CertPath    string `yaml:"cert_path"`
	KeyPath     string `yaml:"key_path"`
	ForceUseTLS bool   `yaml:"force_use_tls"`
}

type ConnectionOptions struct {
	Type string
	Name string

	ConnectionString string
}
