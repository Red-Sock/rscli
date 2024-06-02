package go_actions

import (
	"os"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
)

type BuildProjectAction struct{}

func (a BuildProjectAction) Do(p interfaces.Project) error {
	renamer.ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build()
	if err != nil {
		return err
	}
	return nil
}
func (a BuildProjectAction) NameInAction() string {
	return "Building project"
}

type BuildConfigAction struct{}

func (a BuildConfigAction) Do(p interfaces.Project) error {
	b, err := p.GetConfig().Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshaling config")
	}

	err = os.WriteFile(p.GetConfig().Path, b, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing config to file")
	}

	return nil
}
func (a BuildConfigAction) NameInAction() string {
	return "Building config"
}
