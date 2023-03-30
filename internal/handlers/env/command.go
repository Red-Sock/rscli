package env

import (
	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/handlers/shared"
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	envscripts "github.com/Red-Sock/rscli/plugins/environment/scripts"
)

const Command = "env"

type Handler struct {
	progs map[string]func(args []string) error
}

func NewHandler() *Handler {
	const (
		create = "create"
		setup  = "set-up"
	)

	return &Handler{
		progs: map[string]func(args []string) error{
			create: func(_ []string) error {
				return envscripts.RunCreate()
			},
			setup: func(_ []string) error {
				return envscripts.RunSetUp(nil)
			},
			"help": func(_ []string) error {
				println(shared_ui.Header +
					create + "- create new environment. Run this in root directory where projects are stored\n" +
					setup + " - setting up and update environment for projects\n")
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
