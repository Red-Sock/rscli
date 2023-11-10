package rw

import (
	"bytes"
	"io"
	"sync"
)

type RW struct {
	b []byte
	c int
	l sync.Mutex
}

func (r *RW) Write(b []byte) (int, error) {
	r.l.Lock()
	r.b = append(r.b, b...)
	r.l.Unlock()

	return len(b), nil
}

func (r *RW) WriteByte(b byte) error {
	r.l.Lock()
	r.b = append(r.b, b)
	r.l.Unlock()

	return nil
}

func (r *RW) GetReader() io.Reader {
	r.l.Lock()

	out := bytes.NewReader(r.b)
	r.b = make([]byte, 0, len(r.b))
	r.l.Unlock()
	return out
}

func (r *RW) Bytes() []byte {
	b := make([]byte, len(r.b))
	copy(b, r.b)

	return b
}
