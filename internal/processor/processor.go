package processor

import (
	"os"

	"github.com/Red-Sock/toolbox"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
)

const (
	PathFlag = "path"
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

func (p *Processor) LoadProject(cmd *cobra.Command) (proj *project.Project, err error) {
	var pathToProject string

	if cmd != nil {
		pathToProject = cmd.Flag(PathFlag).Value.String()
	}

	pathToProject = toolbox.Coalesce(pathToProject, p.WD)

	proj, err = project.LoadProject(pathToProject, p.RscliConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error loading project")
	}

	return proj, nil
}
