package config

import (
	"time"

	"github.com/godverv/matreshka"
)

func AppInfo() matreshka.AppInfo {
	return defaultConfig.AppConfig.AppInfo
}

func Api() matreshka.Servers {
	return defaultConfig.AppConfig.Servers
}

func Resources() matreshka.Resources {
	return defaultConfig.AppConfig.Resources
}

func GetInt(key string) (out int) {
	return defaultConfig.GetInt(key)
}
func GetString(key string) (out string) {
	return defaultConfig.GetString(key)
}
func GetBool(key string) (out bool) {
	return defaultConfig.GetBool(key)
}
func GetDuration(key string) (out time.Duration) {
	return defaultConfig.GetDuration(key)
}

func TryGetInt(key string) (out int, err error) {
	return defaultConfig.TryGetInt(key)
}
func TryGetString(key string) (out string, err error) {
	return defaultConfig.TryGetString(key)
}
func TryGetBool(key string) (out bool, err error) {
	return defaultConfig.TryGetBool(key)
}
func TryGetDuration(key string) (t time.Duration, err error) {
	return defaultConfig.TryGetDuration(key)
}
