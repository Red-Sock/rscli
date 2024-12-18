package go_actions

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/bins/makefile"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

const goBin = "go"

type GoFmt struct{}

func (a GoFmt) Do(p project.IProject) error {
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
func (a GoFmt) NameInAction() string {
	return "Performing project fix up"
}

type RunGoTidyAction struct{}

func (a RunGoTidyAction) Do(p project.IProject) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return rerrors.Wrap(err, "error executing go mod tidy")
	}

	err = GoFmt{}.Do(p)
	if err != nil {
		return rerrors.Wrap(err, "error formatting project")
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

func (a RunMakeGenAction) Do(p project.IProject) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	err := makefile.Install()
	if err != nil {
		return rerrors.Wrap(err, "error installing makefile")
	}

	_, err = makefile.Run(p.GetProjectPath(), patterns.RscliMakefileFile, patterns.GenCommand)
	if err != nil {
		return rerrors.Wrap(err, "error running rscli generate command")
	}
	return nil
}
func (a RunMakeGenAction) NameInAction() string {
	return "Running `make gen`"
}

type UpdateAllPackages struct{}

func (a UpdateAllPackages) Do(p project.IProject) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"get", "-u", "all"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a UpdateAllPackages) NameInAction() string {
	return "Updating packages to latest version"
}
