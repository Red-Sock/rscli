package internal

import (
	"github.com/Red-Sock/rscli/internal/handlers/env"
	"github.com/Red-Sock/rscli/internal/handlers/project"
)

type Handle interface {
	Do(args []string) error
}

var h = handler{
	handles: map[string]Handle{
		env.Command:     env.NewHandler(),
		project.Command: project.NewHandler(),
	},
}

type handler struct {
	handles map[string]Handle
}
