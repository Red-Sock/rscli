package env

import (
	"bytes"
	_ "embed"
	"os"
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
)

const (
	equals   = byte('=')
	lineSkip = byte('\n')
)

type Container struct {
	content []Variable
}

type Variable struct {
	Name  string
	Value string
}

func ReadContainer(pth string) (*Container, error) {
	f, err := os.ReadFile(pth)
	if err != nil {
		return nil, err
	}

	return NewEnvContainer(f)
}

func NewEnvContainer(src []byte) (*Container, error) {
	es := &Container{}

	return es, es.UnmarshalEnv(src)
}

func Copy(container *Container) *Container {
	c := &Container{
		content: make([]Variable, len(container.content)),
	}

	copy(c.content, container.content)
	return c
}

func (e *Container) Content() []Variable {
	return e.content
}

func (e *Container) Append(name string, content string) {
	for idx, item := range e.content {
		if item.Name == name {
			e.content[idx].Value = content
			return
		}
	}

	e.content = append(e.content, Variable{Name: name, Value: content})
}

func (e *Container) MarshalEnv() []byte {
	sb := bytes.Buffer{}

	for _, v := range e.content {
		sb.Write([]byte(v.Name))
		sb.WriteByte(equals)
		sb.Write([]byte(v.Value))
		sb.WriteByte(lineSkip)
	}

	return sb.Bytes()
}

func (e *Container) UnmarshalEnv(b []byte) error {
	if b == nil {
		return nil
	}

	splited := bytes.Split(b, []byte{lineSkip})

	e.content = make([]Variable, len(splited))

	for idx, item := range splited {
		line := bytes.Split(item, []byte{equals})
		if len(line) > 0 {
			e.content[idx].Name = string(line[0])
		}
		if len(line) == 2 {
			e.content[idx].Value = string(line[1])
		}
	}

	return nil
}

func (e *Container) RemoveEmpty() {
	newEnvs := make([]Variable, 0, len(e.content)/2)
	for _, item := range e.content {
		if item.Value != "" && item.Name != "" {
			newEnvs = append(newEnvs, item)
		}
	}

	e.content = newEnvs
}

func (e *Container) GetPorts() ([]uint16, error) {
	out := make([]uint16, 0, len(e.content)/2)
	for _, item := range e.content {
		if strings.HasSuffix(item.Name, patterns.PortSuffix) {
			port, err := strconv.ParseUint(item.Value, 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing port for env variable "+item.Name)
			}
			out = append(out, uint16(port))
		}
	}

	return out, nil
}

func (e *Container) Contains(variable Variable) bool {
	for idx := range e.content {
		if e.content[idx].Name == variable.Name && e.content[idx].Value == variable.Value {
			return true
		}
	}

	return false
}
