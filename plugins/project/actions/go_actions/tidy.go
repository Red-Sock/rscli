package go_actions

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/bins/makefile"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

const goBin = "go"

type RunGoFmtAction struct{}

func (a RunGoFmtAction) Do(p proj_interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"fmt", "./..."},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return err
	}

	return nil
}
func (a RunGoFmtAction) NameInAction() string {
	return "Performing project fix up"
}

type RunGoTidyAction struct{}

func (a RunGoTidyAction) Do(p proj_interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error executing go mod tidy")
	}

	err = RunGoFmtAction{}.Do(p)
	if err != nil {
		return errors.Wrap(err, "error formatting project")
	}

	return nil
}
func (a RunGoTidyAction) NameInAction() string {
	return "Cleaning up the project"
}

type RunMakeGenAction struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a RunMakeGenAction) Do(p proj_interfaces.Project) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	err := makefile.Install()
	if err != nil {
		return errors.Wrap(err, "error installing makefile")
	}

	err = makefile.Run(p.GetProjectPath(), projpatterns.Makefile, projpatterns.GenCommand)
	if err != nil {
		return errors.Wrap(err, "error generating")
	}
	return nil
}
func (a RunMakeGenAction) NameInAction() string {
	return "Running `make gen`"
}
