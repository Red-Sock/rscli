package go_actions

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type PrepareProjectStructureAction struct {
}

func (a PrepareProjectStructureAction) Do(p interfaces.Project) error {
	cmd := &folder.Folder{Name: projpatterns.CmdFolder}

	cmd.Inner = append(cmd.Inner, &folder.Folder{
		Name: p.GetShortName(),
		Inner: []*folder.Folder{
			{
				Name:    projpatterns.MainFileName,
				Content: projpatterns.MainFile,
			},
		},
	})

	rootF := p.GetFolder()
	rootF.Add(cmd)

	rootF.Add(&folder.Folder{Name: projpatterns.ConfigsFolder})

	rootF.Add(&folder.Folder{Name: projpatterns.InternalFolder})

	rootF.Add(&folder.Folder{
		Name: projpatterns.PkgFolder,
		Inner: []*folder.Folder{
			{Name: projpatterns.SwaggerFolder},
			{Name: projpatterns.ApiFolder},
		},
	})

	rootF.Add(&folder.Folder{
		Name:    path.Join(projpatterns.InternalFolder, projpatterns.UtilsFolder, projpatterns.CloserFolder, projpatterns.CloserFile),
		Content: projpatterns.UtilsCloserFile,
	})

	return nil
}
func (a PrepareProjectStructureAction) NameInAction() string {
	return "Preparing project structure"
}

type PrepareExamplesFoldersAction struct{}

func (a PrepareExamplesFoldersAction) Do(p interfaces.Project) error {

	if p.GetFolder().GetByPath(projpatterns.ExamplesFolder, projpatterns.ExamplesHttpEnvFile) != nil {
		return nil
	}

	type envs struct {
		Dev       map[string]string `json:"dev"`
		DevDocker map[string]string `json:"dev-docker"`
	}
	var e = envs{
		Dev:       map[string]string{},
		DevDocker: map[string]string{},
	}

	servers, err := p.GetConfig().GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining server options")
	}
	for _, item := range servers {
		portStr := strconv.FormatUint(uint64(item.GetPort()), 10)
		e.Dev[item.GetName()] = "0.0.0.0:" + portStr
		e.DevDocker[item.GetName()] = "0.0.0.0:1" + portStr
	}

	exampleFile, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		return errors.Wrap(err, "error marshalling example file")
	}

	p.GetFolder().Add(&folder.Folder{
		Name: projpatterns.ExamplesFolder,
		Inner: []*folder.Folder{
			{
				Name:    projpatterns.ExampleFileName,
				Content: projpatterns.ApiHTTP,
			},
			{
				Name:    projpatterns.ExamplesHttpEnvFile,
				Content: exampleFile,
			},
		},
	})
	return nil
}
func (a PrepareExamplesFoldersAction) NameInAction() string {
	return "Preparing examples folders"
}

type PrepareEnvironmentFoldersAction struct{}

func (a PrepareEnvironmentFoldersAction) Do(p interfaces.Project) error {
	p.GetFolder().Add(
		[]*folder.Folder{
			{
				Name:    projpatterns.Dockerfile.Name,
				Content: projpatterns.Dockerfile.Content,
			},
			{
				Name:    projpatterns.ReadMeFileName,
				Content: bytes.ReplaceAll(projpatterns.Readme, []byte("{{PROJECT_NAME}}"), []byte(p.GetName())),
			},
			{
				Name:    projpatterns.GitignoreFileName,
				Content: projpatterns.GitIgnore,
			},
			{
				Name:    projpatterns.GolangCIYamlFileName,
				Content: projpatterns.Linter,
			},
			{
				Name:    projpatterns.RsCliMkFileName,
				Content: projpatterns.RscliMK,
			},
		}...,
	)
	return nil
}
func (a PrepareEnvironmentFoldersAction) NameInAction() string {
	return "Preparing environment folder"
}
