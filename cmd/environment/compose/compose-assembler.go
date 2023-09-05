package compose

import (
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/trace-errors"
)

type ComposeAssembler struct {
	Name     string                        `yaml:"-"`
	Services map[string]*ContainerSettings `yaml:"services"`
	Network  map[string]interface{}        `yaml:"networks"`
}

type ContainerSettings struct {
	Image       string            `yaml:"image"`
	WorkDir     string            `yaml:"working_dir,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
}

func NewComposeAssembler(src []byte, pName string) (*ComposeAssembler, error) {
	ca := &ComposeAssembler{
		Name:     pName,
		Services: map[string]*ContainerSettings{},
		Network:  map[string]interface{}{},
	}

	err := yaml.Unmarshal(src, ca)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing src docker-compose file")
	}

	return ca, nil
}

func (c *ComposeAssembler) AppendService(name string, service ContainerSettings) {
	c.Services[name] = &service
}

func (c *ComposeAssembler) Marshal() ([]byte, error) {
	return yaml.Marshal(*c)
}
