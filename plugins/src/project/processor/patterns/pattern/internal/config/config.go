package config

import (
	"flag"
	"github.com/pkg/errors"
	"os"
	"strings"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound    = errors.New("no such key")
	ErrCannotParse = errors.New("couldn't parse value")
)

type Config map[configKey]any

func ReadConfig() (*Config, error) {
	var pth string

	flag.StringVar(&pth, "config", "./config/dev.yaml", "Path to configuration file")
	flag.Parse()

	r, err := os.Open(pth)
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
