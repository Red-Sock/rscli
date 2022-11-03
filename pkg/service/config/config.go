package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type orderedMap map[string]interface{}

type Config struct {
	content         orderedMap
	pth             string
	isForceOverride bool
}

func NewConfig(opts map[string][]string) (*Config, error) {
	cfg, err := buildConfig(opts)
	if err != nil {
		return nil, err
	}

	_, isForceOverride := opts[forceOverride]

	var defaultPath string

	if p, ok := opts[configPath]; ok && len(p) > 0 {
		defaultPath = p[0]
	} else {
		defaultPath = path.Join("./", DefaultDir, FileName)
	}

	return &Config{
		content:         cfg,
		pth:             defaultPath,
		isForceOverride: isForceOverride,
	}, nil
}

func (c *Config) SetFolderPath(pth string) (err error) {
	st, _ := os.Stat(path.Join(pth, FileName))
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

	dir, _ := path.Split(c.pth)

	_ = os.Mkdir(dir, 0775)

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
	out[DataSourceKey] = ds

	serverOptions, err := buildServerOptions(opts)
	if err != nil {
		return nil, err
	}

	out[ServerOptsKey] = serverOptions

	return out, nil
}

func buildDataSources(opts map[string][]string) (map[string]interface{}, error) {
	cfg := make([]map[string]interface{}, 0, len(opts))

	for f, args := range opts {
		switch strings.Replace(f, "-", "", -1) {
		case SourceNamePg:
			cfg = append(cfg, DefaultPgPattern(args))
		case SourceNameRds:
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

func buildServerOptions(opts map[string][]string) (map[string]interface{}, error) {
	cfg := make([]map[string]interface{}, 0, len(opts))

	for f, args := range opts {
		switch strings.Replace(f, "-", "", -1) {
		case RESTHTTPServer:
			cfg = append(cfg, DefaultHTTPPattern(args))
		case GRPCServer:
			cfg = append(cfg, DefaultGRPCPattern(args))
		}
	}
	out := make(map[string]interface{})
	for _, item := range cfg {
		for k, v := range item {
			if _, ok := out[k]; ok {
				return nil, fmt.Errorf("colliding names for transport layers: %s", k)
			}
			out[k] = v
		}
	}
	return out, nil
}
