package env

import (
	"bytes"
	_ "embed"
	"os"
	"strconv"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/envpatterns"
)

const (
	equals   = byte('=')
	lineSkip = byte('\n')
)

type Container struct {
	Content []Variable
}

type Variable struct {
	Name  string
	Value string
}

func ReadContainer(pth string) (*Container, error) {
	f, err := os.ReadFile(pth)
	if err != nil {
		if !rerrors.Is(err, os.ErrNotExist) {
			return nil, rerrors.Wrap(err, "error reading container")
		}
	}

	if len(f) != 0 {
		return NewEnvContainer(f)
	}

	return NewEnvContainer(nil)
}

func NewEnvContainer(src []byte) (*Container, error) {
	es := &Container{}
	if len(src) == 0 {
		return es, nil
	}

	return es, es.UnmarshalEnv(src)
}

func (e *Container) GetContent() []Variable {
	return e.Content
}

func (e *Container) Append(v ...Variable) {
	for _, item := range v {
		e.append(item)
	}
}

func (e *Container) AppendRaw(name string, content string) {
	for idx, item := range e.Content {
		if item.Name == name {
			e.Content[idx].Value = content
			return
		}
	}

	e.Content = append(e.Content, Variable{Name: name, Value: content})
}

func (e *Container) MarshalEnv() []byte {
	if len(e.Content) == 0 {
		return []byte{}
	}
	sb := bytes.Buffer{}

	for _, v := range e.Content {
		sb.Write([]byte(v.Name))
		if v.Name != "" && v.Name[0] != '#' {
			sb.WriteByte(equals)
		}

		sb.Write([]byte(v.Value))
		sb.WriteByte(lineSkip)
	}
	out := sb.Bytes()
	return out[:len(out)-1]
}

func (e *Container) UnmarshalEnv(b []byte) error {
	if b == nil {
		return nil
	}

	splited := bytes.Split(b, []byte{lineSkip})

	e.Content = make([]Variable, len(splited))

	for idx, item := range splited {
		line := bytes.Split(item, []byte{equals})
		if len(line) > 0 {
			e.Content[idx].Name = string(line[0])
		}
		if len(line) == 2 {
			e.Content[idx].Value = string(line[1])
		}
	}

	return nil
}

type IntVariable struct {
	Name  string
	Value uint16
}

func (e *Container) GetPortValues() ([]IntVariable, error) {
	out := make([]IntVariable, 0, len(e.Content)/2)
	for _, item := range e.Content {
		if strings.HasSuffix(item.Name, envpatterns.PortSuffix) {
			port, err := strconv.ParseUint(item.Value, 10, 16)
			if err == nil {
				out = append(out, IntVariable{
					Name:  item.Name,
					Value: uint16(port),
				})
			}
		}
	}

	return out, nil
}

func (e *Container) Contains(variable Variable) bool {
	for idx := range e.Content {
		if e.Content[idx].Name == variable.Name && e.Content[idx].Value == variable.Value {
			return true
		}
	}

	return false
}

func (e *Container) ContainsByName(name string) bool {
	for _, item := range e.Content {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (e *Container) Rename(oldName, newName string) {
	for idx := range e.Content {
		if e.Content[idx].Name == oldName {
			e.Content[idx].Name = newName
			return
		}
	}
}

func (e *Container) Remove(name string) {
	for idx, envVar := range e.Content {
		if envVar.Name == name {
			e.Content[0], e.Content[idx] = e.Content[idx], e.Content[0]
			break
		}
	}
	e.Content = e.Content[1:]
}

func (e *Container) append(v Variable) {
	for idx, item := range e.Content {
		if item.Name == v.Name {
			e.Content[idx].Value = v.Value
			return
		}
	}

	e.Content = append(e.Content, Variable{Name: v.Name, Value: v.Value})
}
