package env

import (
	"github.com/Red-Sock/rscli/internal/handlers/shared"
	envscripts "github.com/Red-Sock/rscli/plugins/environment/scripts"
	"github.com/pkg/errors"
)

const Command = "env"

type Handler struct {
	progs map[string]func(args []string) error
}

func NewHandler() *Handler {
	return &Handler{
		progs: map[string]func(args []string) error{
			"create": func(_ []string) error {
				return envscripts.RunCreate()
			},
		},
	}
}

func (h *Handler) Do(args []string) error {
	if len(args) == 0 {
		return shared.ErrNoArguments
	}

	hl, ok := h.progs[args[0]]
	if !ok {
		return errors.Wrapf(shared.ErrUnknownHandler, "creating %s is out of my abilities", args[0])
	}

	return hl(args[1:])
}
