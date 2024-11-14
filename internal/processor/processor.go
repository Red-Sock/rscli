package processor

import (
	"os"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
)

// Processor - represents a single process of execution.
// e.g. rscli project tidy - calls a cmd/project/tidy Processor and executes it
// Contains all basic necessary information and primitives for CLI utility
type Processor struct {
	IO          io.IO
	RscliConfig *config.RsCliConfig
	WD          string
}

type opt func(p *Processor)

func New(opts ...opt) Processor {
	p := Processor{}

	for _, o := range opts {
		o(&p)
	}

	if p.IO == nil {
		p.IO = io.StdIO{}
	}

	if p.WD == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		p.WD = wd
	}

	if p.RscliConfig == nil {
		p.RscliConfig = config.GetConfig()
	}

	return p
}
