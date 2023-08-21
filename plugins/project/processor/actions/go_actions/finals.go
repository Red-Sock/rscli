package go_actions

import (
	"os"
	"path"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/go_actions/tidy"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
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

type FixupProjectAction struct{}

func (a FixupProjectAction) Do(p interfaces.Project) error {
	wd, _ := os.Getwd()
	wd = path.Join(wd, p.GetName())

	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: wd,
	})
	if err != nil {
		return errors.Wrap(err, "error executing go mod tidy")
	}

	_, err = cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"fmt", "./..."},
		WorkDir: wd,
	})
	if err != nil {
		return err
	}

	return nil
}
func (a FixupProjectAction) NameInAction() string {
	return "Performing project fix up"
}

type TidyAction struct{}

func (a TidyAction) Do(p interfaces.Project) error {
	err := tidy.Api(p)
	if err != nil {
		return errors.Wrap(err, "error during api tiding")
	}

	err = tidy.Config(p)
	if err != nil {
		return errors.Wrap(err, "error during config tiding")
	}

	err = tidy.DataSources(p)
	if err != nil {
		return errors.Wrap(err, "error during data source tiding")
	}

	ReplaceProjectName(p.GetName(), p.GetFolder())

	err = BuildGoConfigFolderAction{}.Do(p)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	return errors.Wrap(p.GetFolder().Build(), "error building project")
}
func (a TidyAction) NameInAction() string {
	return "Cleaning up the project"
}

// helping functions
