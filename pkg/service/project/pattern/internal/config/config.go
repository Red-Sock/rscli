package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config map[string]string

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
	return nil, nil
}
