package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/helpers/cases"
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
		variables[cases.SnakeToCamel(v)] = v
	}

	return variables, nil
}

type cfgKeysBuilder map[string]interface{}

func (c *cfgKeysBuilder) extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(cfgKeysBuilder); ok {
			values, err := c.extractVariables(prefix+"_"+k, newMap)
			if err != nil {
				return nil, err
			}
			out = append(out, values...)
		} else {
			k = prefix + "_" + k

			out = append(out, k[1:])
		}
	}
	return out, nil
}
