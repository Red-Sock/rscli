package go_project

import (
	"os"
	"strings"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

type Validator func(p project.Project) error

type Project struct {
	Name        string
	ProjectPath string

	Cfg *config.Config

	Actions []actions.Action

	projType project.Type

	validators []Validator
	root       folder.Folder
}

func (p *Project) GetShortName() string {
	name := p.Name

	if idx := strings.LastIndex(name, string(os.PathSeparator)); idx != -1 {
		name = name[idx+1:]
	}

	return name
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetFolder() *folder.Folder {
	return &p.root
}

func (p *Project) GetConfig() *config.Config {
	return p.Cfg
}

func (p *Project) GetProjectPath() string {
	return p.ProjectPath
}

func (p *Project) GetActionNames() []string {
	out := make([]string, 0, len(p.Actions))
	for _, a := range p.Actions {
		out = append(out, a.NameInAction())
	}
	return out
}

func (p *Project) Build() (<-chan string, <-chan error) {
	progressCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(progressCh)
		defer close(errCh)

		for _, a := range p.Actions {
			progressCh <- a.NameInAction()
			if err := a.Do(p); err != nil {
				errCh <- err
				return
			}
		}

	}()

	return progressCh, errCh
}

func (p *Project) Validate() error {
	var errs []error
	for _, v := range p.validators {
		if err := v(p); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	globalErr := errors.New("error while validating the project")
	for _, e := range errs {
		globalErr = errors.Wrap(globalErr, e.Error())
	}

	return globalErr
}

func (p *Project) GetType() project.Type {
	return p.projType
}
