package help

import (
	"github.com/Red-Sock/rscli/internal/shared-ui"

	"github.com/Red-Sock/rscli/internal/handlers/env"
	"github.com/Red-Sock/rscli/internal/handlers/project"
	"github.com/Red-Sock/rscli/internal/handlers/version"
)

const Command = "help"

type Handler struct{}

func (h *Handler) Do(args []string) error {
	println(shared_ui.Header +
		env.Command + " - local dev environment setup\n" +
		project.Command + " - manage project setting or create new project\n" +
		version.Command + " - get current version of RSCLI\n")
	return nil
}
