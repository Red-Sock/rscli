package patterns

import (
	"bytes"
	_ "embed"
)

const (
	equals   = byte('=')
	lineSkip = byte('\n')
)

type EnvService struct {
	content []EnvironmentValue
}

type EnvironmentValue struct {
	Name  string
	Value string
}

func NewEnvService(src []byte) (*EnvService, error) {
	es := &EnvService{}

	return es, es.UnmarshalEnv(src)
}

func (e *EnvService) Content() []EnvironmentValue {
	return e.content
}

func (e *EnvService) Append(name string, content string) {
	e.content = append(e.content, EnvironmentValue{Name: name, Value: content})
}

func (e *EnvService) MarshalEnv() []byte {
	sb := bytes.Buffer{}

	for _, v := range e.content {
		sb.Write([]byte(v.Name))
		sb.WriteByte(equals)
		sb.Write([]byte(v.Value))
		sb.WriteByte(lineSkip)
	}

	return sb.Bytes()
}

func (e *EnvService) UnmarshalEnv(b []byte) error {
	if b == nil {
		return nil
	}

	splited := bytes.Split(b, []byte{lineSkip})

	e.content = make([]EnvironmentValue, len(splited))

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
