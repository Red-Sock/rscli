package config

import (
	"os"

	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"gopkg.in/yaml.v3"
)

type Config struct {
	*matreshka.AppConfig

	pth string `yaml:"-"`
}

func ReadConfig(pth string) (*Config, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}

	c := &Config{
		pth: pth,
	}
	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding config to struct")
	}

	return c, nil
}

func (c *Config) GetPath() string {
	return c.pth
}
