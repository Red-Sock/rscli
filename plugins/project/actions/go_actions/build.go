package go_actions

import (
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
)

type BuildProjectAction struct{}

func (a BuildProjectAction) Do(p project.IProject) error {
	renamer.ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build(p.GetProjectPath())
	if err != nil {
		return err
	}
	return nil
}
func (a BuildProjectAction) NameInAction() string {
	return "Building project"
}
