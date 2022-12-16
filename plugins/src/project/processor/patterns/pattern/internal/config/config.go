package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type Config map[configKey]any

func ReadConfig() (*Config, error) {
	var pth string

	flag.StringVar(&pth, "config", "./config/dev.yaml", "Path to configuration file")
	flag.Parse()

	r, err := os.Open(pth)
	if err != nil {
		return nil, fmt.Errorf("os.Open %w", err)
	}

	var cfg map[string]interface{}
	err = yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("yaml.Decode %w", err)
	}

	res, err := extractVariables("", cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config %w", err)
	}

	c := Config(res)

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
		return "", errors.New(fmt.Sprintf("no string value for key %s was found", key))
	}

	out, ok := v.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("couldn't parse value %v to string", v))
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
		return 0, errors.New(fmt.Sprintf("no string value for key %s was found", key))
	}

	out, ok := v.(int)
	if !ok {
		return 0, errors.New(fmt.Sprintf("couldn't parse value %v to string", v))
	}
	return out, nil
}

func extractVariables(prefix string, in map[string]interface{}) (out map[configKey]any, err error) {
	out = make(map[configKey]any)
	for k, v := range in {
		if newMap, ok := v.(map[string]interface{}); ok {
			values, err := extractVariables(mergeParts("_", prefix, k), newMap)
			if err != nil {
				return nil, err
			}
			maps.Copy(out, values)
		} else {
			out[configKey(mergeParts("_", prefix, k))] = v
		}
	}
	return out, nil
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
