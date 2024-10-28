package actions

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
)

type Action interface {
	Do(p project.IProject) error
	NameInAction() string
}

type ActionPerformer struct {
	printer io.IO
	proj    project.IProject
}

func NewActionPerformer(printer io.IO, proj project.IProject) ActionPerformer {
	return ActionPerformer{
		printer: printer,
		proj:    proj,
	}
}

func (a *ActionPerformer) Tidy() error {
	acts := GetTidyActionsForProject(a.proj.GetType())

	for _, ac := range acts {
		a.printer.Println(ac.NameInAction())

		err := ac.Do(a.proj)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}
