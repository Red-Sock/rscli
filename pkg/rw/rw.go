package rw

import (
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

func (r *RW) Read(b []byte) (int, error) {
	r.l.Lock()

	n := cap(b)
	if cap(b) > len(r.b) {
		n = len(r.b)
	}
	for i := 0; i < n; i++ {
		b[i] = r.b[i]
	}
	var err error
	if len(b) >= len(r.b) {
		r.b = make([]byte, 0, len(b))
		err = io.EOF
	} else {
		r.b = r.b[len(b):]
	}

	r.l.Unlock()

	return len(b), err
}
