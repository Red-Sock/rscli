package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func KeysFromConfig(pathToConfig string) (map[string]string, error) {
	cfgBytes, err := os.ReadFile(pathToConfig)
	if err != nil {
		return nil, err
	}

	cfg := make(cfgKeysBuilder)
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, err
	}

	vars, err := cfg.extractVariables("", cfg)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]string, len(cfg))
	for _, v := range vars {
		parts := strings.Split(v[1:], "_")
		for i := range parts {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
		variables[strings.Join(parts, "")] = v[1:]
	}

	return variables, nil
}

type cfgKeysBuilder map[string]interface{}

func (c *cfgKeysBuilder) extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(map[string]interface{}); ok {
			values, err := c.extractVariables(prefix+"_"+k, newMap)
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
