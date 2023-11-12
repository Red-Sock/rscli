package go_actions

import (
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
)

type BuildProjectAction struct{}

func (a BuildProjectAction) Do(p interfaces.Project) error {
	ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build()
	if err != nil {
		return err
	}
	return nil
}
func (a BuildProjectAction) NameInAction() string {
	return "Building project"
}
