package config

type Config struct {
	AppInfo     AppInfo                `yaml:"app_info"`
	Server      map[string]interface{} `yaml:"server"`
	DataSources map[string]interface{} `yaml:"data_sources"`
}

type AppInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]interface{}{},
		DataSources: map[string]interface{}{},
	}
}
