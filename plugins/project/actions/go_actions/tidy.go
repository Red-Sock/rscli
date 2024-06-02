package go_actions

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/bins/makefile"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

const goBin = "go"

type InitGoModAction struct{}

func (a InitGoModAction) Do(p interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "init", p.GetName()},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error executing go mod init")
	}

	goMod, err := os.OpenFile(path.Join(p.GetProjectPath(), "go.mod"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer func() {
		err2 := goMod.Close()
		if err2 != nil {
			if err != nil {
				err = errors.Wrap(err, "error on closing"+err2.Error())
			} else {
				err = err2
			}
		}
	}()

	return nil
}
func (a InitGoModAction) NameInAction() string {
	return "Initiating go project"
}

type RunGoFmtAction struct{}

func (a RunGoFmtAction) Do(p interfaces.Project) error {
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

func (a RunGoTidyAction) Do(p interfaces.Project) error {
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

func (a RunMakeGenAction) Do(p interfaces.Project) error {
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
