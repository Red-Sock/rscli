package project

import (
	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/handlers/shared"
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
)

const Command = "project"

type Handler struct {
	progs map[string]func(args []string) error
}

func NewHandler() *Handler {
	const (
		create = "create"
		tidy   = "tidy"
	)
	return &Handler{
		progs: map[string]func(args []string) error{
			create: createProject,
			tidy:   tidyProject,
			"help": func(_ []string) error {
				println(shared_ui.Header +
					create + " - creates new project\n" +
					tidy + " - clears project and updates it to a newer version of RSCLI\n")
				return nil
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
