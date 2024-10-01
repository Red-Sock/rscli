package go_actions

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators/app_struct_generators"
)

type InitGoModAction struct{}

func (a InitGoModAction) Do(p project.Project) error {
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

type InitGoProjectApp struct{}

func (a InitGoProjectApp) Do(p project.Project) error {
	appFolderPath := path.Join(projpatterns.InternalFolder, projpatterns.AppFolder)
	appFolder := p.GetFolder().GetByPath(appFolderPath)
	if appFolder == nil {
		appFolder = &folder.Folder{
			Name: path.Join(projpatterns.InternalFolder, projpatterns.AppFolder),
		}
		p.GetFolder().Add(appFolder)
	}

	appFiles, err := app_struct_generators.GenerateAppFiles(p)
	if err != nil {
		return errors.Wrap(err, "error generating app file")
	}

	for fileName, fileContent := range appFiles {
		appFolder.Add(
			&folder.Folder{
				Name:    fileName,
				Content: fileContent,
			})
	}

	return nil
}

func (a InitGoProjectApp) NameInAction() string {
	return "Generating app skeleton"
}
