package config

type Config struct {
	AppInfo     AppInfo                `yaml:"app_info,omitempty"`
	Server      map[string]interface{} `yaml:"server"`
	DataSources map[string]interface{} `yaml:"data_sources"`
}

type AppInfo struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]interface{}{},
		DataSources: map[string]interface{}{},
	}
}
