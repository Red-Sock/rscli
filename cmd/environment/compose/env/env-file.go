package env

import (
	"bytes"
	_ "embed"
	"os"
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
