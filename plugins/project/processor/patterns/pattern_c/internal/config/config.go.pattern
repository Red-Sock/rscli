package config

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

const (
	devConfigPath  = "./config/dev.yaml"
	prodConfigPath = "./config/config.yaml"
)

var (
	ErrNotFound    = errors.New("no such key")
	ErrCannotParse = errors.New("couldn't parse value")
)

type Config map[configKey]any

func ReadConfig() (*Config, error) {
	var cfgPath string
	var isDevBuild bool

	flag.StringVar(&cfgPath, "config", "", "Path to configuration file")
	flag.BoolVar(&isDevBuild, "dev", false, "Flag turns on a dev config at ./config/dev.yaml")
	flag.Parse()

	if cfgPath == "" {
		if isDevBuild {
			cfgPath = devConfigPath
		} else {
			cfgPath = prodConfigPath
		}
	}

	r, err := os.Open(cfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}

	var cfg map[string]interface{}

	err = yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "yaml.Decode")
	}

	c := Config(extractVariables("", cfg))

	return &c, nil
}

func (c Config) GetString(key configKey) string {
	v := c[key]
	if v == nil {
		return ""
	}

	out, ok := v.(string)
	if !ok {
		return ""
	}

	return out
}
func (c Config) TryGetString(key configKey) (string, error) {
	v := c[key]
	if v == nil {
		return "", errors.Wrapf(ErrNotFound, "absent key: %s", key)
	}

	out, ok := v.(string)
	if !ok {
		return "", errors.Wrapf(ErrCannotParse, "values is: %v", v)
	}

	return out, nil
}

func (c Config) GetInt(key configKey) int {
	v := c[key]
	if v == nil {
		return 0
	}

	out, ok := v.(int)
	if !ok {
		return 0
	}

	return out
}
func (c Config) TryGetInt(key configKey) (int, error) {
	v := c[key]
	if v == nil {
		return 0, errors.Wrapf(ErrNotFound, "absent key: %s", key)
	}

	out, ok := v.(int)
	if !ok {
		return 0, errors.Wrapf(ErrCannotParse, "value is: %v", v)
	}

	return out, nil
}

func (c Config) GetDuration(key configKey) (out time.Duration, err error) {
	v, ok := c[key]
	if !ok {
		return 0, ErrNotFound
	}

	r, ok := c[key].(string)
	if !ok {
		return 0, errors.Wrapf(ErrCannotParse, "%v of type %T to string", v, v)
	}

	return time.ParseDuration(r)

}

func extractVariables(prefix string, in map[string]interface{}) (out map[configKey]any) {
	out = make(map[configKey]any)

	for k, v := range in {
		if newMap, ok := v.(map[string]interface{}); ok {
			maps.Copy(out, extractVariables(mergeParts("_", prefix, k), newMap))
		} else {
			out[configKey(mergeParts("_", prefix, k))] = v
		}
	}

	return out
}

func mergeParts(delimiter string, parts ...string) string {
	notEmptyParts := make([]string, 0, len(parts))

	for _, p := range parts {
		if p == "" {
			continue
		}

		notEmptyParts = append(notEmptyParts, p)
	}

	return strings.Join(notEmptyParts, delimiter)
}
