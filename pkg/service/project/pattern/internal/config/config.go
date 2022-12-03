package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config map[string]any

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

	variables := make(map[string]any, len(cfg))
	for _, v := range res {
		parts := strings.Split(v[1:], "_")
		for i := range parts {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
		variables[strings.Join(parts, "")] = v[1:]
	}

	c := Config(variables)

	return &c, nil
}

func (c *Config) GetString(key configKey) {

}

func extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(map[string]interface{}); ok {
			values, err := extractVariables(prefix+"_"+k, newMap)
			if err != nil {
				return nil, err
			}
			out = append(out, values...)
		} else {
			out = append(out, prefix+"_"+k)
		}
	}
	return out, nil
}
