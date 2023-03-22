package project

import (
	"github.com/Red-Sock/rscli/internal/handlers/shared"
	"github.com/pkg/errors"
)

const Command = "project"

type Handler struct {
	progs map[string]func(args []string) error
}

func NewHandler() *Handler {
	return &Handler{
		progs: map[string]func(args []string) error{
			"create": createProject,
			"tidy":   tidyProject,
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
