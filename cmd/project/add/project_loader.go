package project

import (
	"github.com/Red-Sock/toolbox"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/plugins/project"
)

func (p *Proc) loadProject(cmd *cobra.Command) (proj *project.Project, err error) {
	var pathToProject string

	if cmd != nil {
		pathToProject = cmd.Flag(pathFlag).Value.String()
	}

	pathToProject = toolbox.Coalesce(pathToProject, p.WD)

	proj, err = project.LoadProject(pathToProject, p.RscliConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error loading project")
	}

	return proj, nil
}
