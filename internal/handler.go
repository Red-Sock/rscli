package internal

import (
	"github.com/Red-Sock/rscli/internal/handlers/env"
	"github.com/Red-Sock/rscli/internal/handlers/help"
	"github.com/Red-Sock/rscli/internal/handlers/project"
	"github.com/Red-Sock/rscli/internal/handlers/version"
)

type Handle interface {
	Do(args []string) error
}

var h = handler{
	handles: map[string]Handle{
		env.Command:     env.NewHandler(),
		project.Command: project.NewHandler(),

		version.Command: &version.Handler{},

		help.Command: &help.Handler{},
	},
}

type handler struct {
	handles map[string]Handle
}
