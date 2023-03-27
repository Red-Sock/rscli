package closer

import (
	"errors"
)

type Closable func() error

var funcs []Closable

func Add(f Closable) {
	funcs = append(funcs, f)
}

func Close() (err error) {
	for _, f := range funcs {
		fErr := f()
		if err != nil {
			err = errors.Join(err, fErr)
		}
	}

	return err
}
