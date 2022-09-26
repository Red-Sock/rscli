package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// root keys names
const (
	dataSource = "data_sources"
)

const (
	sourceNamePg  = "pg"
	sourceNameRds = "rds"

	forceOverride = "fo"
	configPath    = "path"
)

type Config struct {
	content         map[string]interface{}
	pth             string
	isForceOverride bool
}

func NewConfig(opts map[string][]string) (*Config, error) {
	cfg, err := buildConfig(opts)
	if err != nil {
		return nil, err
	}

	_, isForceOverride := opts[forceOverride]

	defaultPath := "./config/local_config.yaml"
	if p, ok := opts[configPath]; ok && len(p) > 0 {
		defaultPath = p[0]
	}

	return &Config{
		content:         cfg,
		pth:             defaultPath,
		isForceOverride: isForceOverride,
	}, nil
}

func (c *Config) SetPath(pth string) (err error) {
	st, _ := os.Stat(pth)
	if st != nil {
		return os.ErrExist
	}

	c.pth = pth
	return nil
}

func (c *Config) GetPath() string {
	return c.pth
}

func (c *Config) TryWrite() (err error) {

	var f *os.File

	_ = os.Mkdir("./config", 0775)

	st, _ := os.Stat(c.pth)
	if st != nil {
		return os.ErrExist
	}
	if f, err = os.Create(c.pth); err != nil {
		return err
	}
	defer f.Close()

	if err = yaml.NewEncoder(f).Encode(c.content); err != nil {
		return err
	}

	return nil
}

func (c *Config) ForceWrite() (err error) {
	_ = os.RemoveAll(c.pth)
	w, err := os.Create(c.pth)
	if err != nil {
		return err
	}
	defer w.Close()

	if err = yaml.NewEncoder(w).Encode(c.content); err != nil {
		return err
	}
	return nil
}

func buildConfig(opts map[string][]string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	ds, err := buildDataSources(opts)
	if err != nil {
		return nil, err
	}
	out[dataSource] = ds
	return out, nil
}

func buildDataSources(opts map[string][]string) (map[string]interface{}, error) {
	cfg := make([]map[string]interface{}, 0, len(opts))

	for f, args := range opts {
		switch strings.Replace(f, "-", "", -1) {
		case sourceNamePg:
			cfg = append(cfg, DefaultPgPattern(args))
		case sourceNameRds:
			cfg = append(cfg, DefaultRdsPattern(args))
		}
	}
	out := make(map[string]interface{})
	for _, item := range cfg {
		for k, v := range item {
			if _, ok := out[k]; ok {
				return nil, fmt.Errorf("colliding names for data sources: %s", k)
			}
			out[k] = v
		}
	}
	return out, nil
}
