package processor

import (
	"os"
	"strings"

	"github.com/go-faster/errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

type Validator func(p interfaces.Project) error

type Project struct {
	Name        string
	ProjectPath string
	Cfg         interfaces.ProjectConfig

	Actions []actions.Action

	RsCLIVersion interfaces.Version

	validators []Validator
	root       folder.Folder
}

func (p *Project) GetShortName() string {
	name := p.Name
	name = name[strings.LastIndex(name, string(os.PathSeparator))+1:]
	return name
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetFolder() *folder.Folder {
	return &p.root
}

func (p *Project) GetConfig() interfaces.ProjectConfig {
	return p.Cfg
}

func (p *Project) GetProjectPath() string {
	return p.ProjectPath
}

func (p *Project) Build() (<-chan string, <-chan error) {
	progressCh := make(chan string, len(p.Actions))
	errCh := make(chan error)

	go func() {
		for _, a := range p.Actions {
			progressCh <- a.NameInAction()
			if err := a.Do(p); err != nil {
				close(progressCh)
				errCh <- err
				close(errCh)
				return
			}
		}
		close(progressCh)
		close(errCh)
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

func (p *Project) GetVersion() interfaces.Version {
	return p.RsCLIVersion
}

func (p *Project) SetVersion(version interfaces.Version) {
	p.RsCLIVersion = version
}
