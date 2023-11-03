package go_actions

import (
	"os"
	"path"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
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

type FormatAction struct{}

func (a FormatAction) Do(p interfaces.Project) error {
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
func (a FormatAction) NameInAction() string {
	return "Performing project fix up"
}

type TidyAction struct{}

func (a TidyAction) Do(p interfaces.Project) error {
	ReplaceProjectName(p.GetName(), p.GetFolder())

	err := PrepareGoConfigFolderAction{}.Do(p)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	err = p.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building project")
	}
	return nil
}
func (a TidyAction) NameInAction() string {
	return "Cleaning up the project"
}
