package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App app
}

type app struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Port    string `yaml:"port"`
}

func ReadConfig() (*Config, error) {
	var pth string

	flag.StringVar(&pth, "config", "./configs/dev.yaml", "Path to configuration file")
	flag.Parse()

	r, err := os.Open(pth)
	if err != nil {
		return nil, fmt.Errorf("os.Open %w", err)
	}

	var cfg Config
	err = yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("yaml.Decode %w", err)
	}
	return &cfg, nil
}
