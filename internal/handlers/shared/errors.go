package shared

import "github.com/pkg/errors"

var (
	ErrNoArguments    = errors.New("create ... what? specify what to create!")
	ErrUnknownHandler = errors.New("sorry, I don't understand.")
)
