package errors

import (
	"fmt"
	"runtime"
)

func Wrap(err error, msg string) error {
	var frames [3]uintptr
	runtime.Callers(2, frames[:])

	return Error{
		innerError: err,
		msg:        msg,
		trace:      frames,
	}
}

func Wrapf(err error, msg string, args ...interface{}) error {
	var frames [3]uintptr
	runtime.Callers(2, frames[:])

	return Error{
		innerError: err,
		msg:        fmt.Sprintf(msg, args...),
		trace:      frames,
	}
}
