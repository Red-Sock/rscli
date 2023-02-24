package create

import (
	envscripts "github.com/Red-Sock/rscli/plugins/environment/scripts"
	"github.com/pkg/errors"
	"os"
)

var (
	ErrNoArguments    = errors.New("create ... what? specify what to create!")
	ErrUnknownHandler = errors.New("sorry, I don't understand.")
)

const Command = "create"

const (
	help = "help"
)

type Handler struct {
	progs map[string]func(args []string) error
}

func NewHandler() *Handler {
	return &Handler{
		progs: map[string]func(args []string) error{
			"env": func(_ []string) error {
				return envscripts.RunCreate()
			},
			"project": createProject,
		},
	}
}

func (h *Handler) Do(args []string) error {
	if len(args) == 0 {
		return ErrNoArguments
	}

	hl, ok := h.progs[args[0]]
	if !ok {
		return errors.Wrapf(ErrUnknownHandler, "creating %s is out of my abilities", os.Args[0])
	}

	return hl(args[1:])
}
