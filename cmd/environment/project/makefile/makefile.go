package makefile

import (
	"bytes"
	"os"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
)

type command struct {
	name     []byte
	phony    []byte
	commands [][]byte
}

type Makefile struct {
	Variables *env.Container
	Commands  []command
}

func ReadMakeFile(pth string) (*Makefile, error) {
	makeFile, err := os.ReadFile(pth)
	if err != nil {
		return nil, errors.Wrap(err, "error reading makefile")
	}

	return NewMakeFile(makeFile)
}

func NewMakeFile(in []byte) (*Makefile, error) {
	lines := bytes.Split(in, []byte{'\n'})
	m := &Makefile{
		Variables: &env.Container{},
	}

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		if index := bytes.Index(l, []byte{'='}); index != -1 {
			m.Variables.Append(env.Variable{
				Name:  string(l[:index]),
				Value: string(l[index+1:]),
			})
		}

		if index := bytes.Index(l, []byte{':'}); index != -1 {
			m.Variables.Append(env.Variable{
				Name:  string(l[:index]),
				Value: string(l[index+1:]),
			})
		}
	}
	return m, nil
}
