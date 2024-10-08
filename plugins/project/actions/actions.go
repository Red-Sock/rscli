package actions

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
)

type Action interface {
	Do(p proj_interfaces.Project) error
	NameInAction() string
}

type ActionPerformer struct {
	printer io.IO
	proj    proj_interfaces.Project
}

func NewActionPerformer(printer io.IO, proj proj_interfaces.Project) ActionPerformer {
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
