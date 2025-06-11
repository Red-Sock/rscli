package compose

import (
	"gopkg.in/yaml.v3"
)

type Compose struct {
	Services map[string]*ContainerSettings `yaml:"services"`
	Network  map[string]interface{}        `yaml:"networks"`
}

type ContainerSettings struct {
	Image       string            `yaml:"image"`
	WorkDir     string            `yaml:"working_dir,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
}

func (c *Compose) AppendService(name string, service ContainerSettings) {
	c.Services[name] = &service
}

func (c *Compose) Marshal() ([]byte, error) {
	return yaml.Marshal(*c)
}
