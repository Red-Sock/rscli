package env

import (
	"github.com/Red-Sock/rscli/internal/handlers/shared"
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	envscripts "github.com/Red-Sock/rscli/plugins/environment/scripts"
)

const Command = "env"

const (
	setup = "set-up"
)

type Handler struct {
	output chan string
	input  chan string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Do(args []string) error {
	if len(args) == 0 {
		return shared.ErrNoArguments
	}

	switch args[0] {
	case setup:
		return h.Setup()
	default:

		println(shared_ui.Header +
			setup + " - setting up and update environment for projects\n")
		return nil

	}
}

func (h *Handler) Setup() error {
	return envscripts.RunSetUp(nil)
}
