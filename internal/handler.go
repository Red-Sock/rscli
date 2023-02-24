package internal

import (
	"github.com/Red-Sock/rscli/internal/handlers/create"
)

type Handle interface {
	Do(args []string) error
}

var h = handler{
	handles: map[string]Handle{
		create.Command: create.NewHandler(),
	},
}

type handler struct {
	handles map[string]Handle
}
