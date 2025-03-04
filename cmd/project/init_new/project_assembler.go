package init_new

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
)

func (p *Proc) createProject(args project.CreateArgs) (project.IProject, error) {
	proj, err := project.CreateProject(args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error during project creation")
	}

	p.IO.Println("Starting project constructor")

	initActions := actions.InitProject(project.TypeGo)
	for _, act := range initActions {
		err = act.Do(proj)
		if err != nil {
			return nil, rerrors.Wrap(err, "error performing init actions")
		}
	}

	p.IO.Println("Project actions performed")

	return proj, nil
}
