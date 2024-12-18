package actions

//go:generate minimock -i ActionPerformer -o ./../../../tests/mocks -g -s "_mock.go"

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/plugins/project"
)

type Action interface {
	Do(p project.IProject) error
	NameInAction() string
}

type ActionPerformer interface {
	Tidy(proj project.IProject) error
}

type actionPerformer struct {
	printer io.IO
}

func NewActionPerformer(printer io.IO) ActionPerformer {
	return &actionPerformer{
		printer: printer,
	}
}

func (a *actionPerformer) Tidy(proj project.IProject) error {
	acts := GetTidyActionsForProject(proj.GetType())

	for _, ac := range acts {
		a.printer.Println(ac.NameInAction())

		err := ac.Do(proj)
		if err != nil {
			return rerrors.Wrap(err)
		}
	}

	return nil
}
