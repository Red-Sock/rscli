package cmd

import (
	"github.com/pkg/errors"
	"io"
	"os/exec"
)

type Request struct {
	Tool    string
	Args    []string
	WorkDir string
}

func Execute(r Request) (message string, err error) {
	cmd := exec.Command(r.Tool, r.Args...)
	if r.WorkDir != "" {
		cmd.Dir = r.WorkDir
	}
	errRW := &RW{}
	cmd.Stderr = errRW

	msgRW := &RW{}
	cmd.Stdout = msgRW

	err = cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, errRW.String())
	}

	return msgRW.String(), err
}

type RW struct {
	b []byte
}

func (r *RW) Write(p []byte) (n int, err error) {
	r.b = append(r.b, p...)
	return len(p), nil
}

func (r *RW) Read(b []byte) (n int, err error) {
	for idx := range b {
		if idx >= len(r.b) {
			break
		}
		b[idx] = r.b[idx]
		n++
	}
	if n == 0 {
		return 0, io.EOF
	}

	r.b = r.b[n:]
	return n, nil
}

func (r *RW) String() string {
	bts, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error parsing message from execution error").Error()
	}
	return string(bts)
}
