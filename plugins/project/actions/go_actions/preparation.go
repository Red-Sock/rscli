package go_actions

import (
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
			projpatterns.MainFile,
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

	rootF.Add(
		projpatterns.UtilsCloserFile.CopyWithNewName(
			path.Join(projpatterns.InternalFolder, projpatterns.UtilsFolder, projpatterns.CloserFolder, projpatterns.UtilsCloserFile.Name)),
	)

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

	for _, item := range p.GetConfig().Server {
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
			projpatterns.ApiHTTP.Copy(),
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
			projpatterns.Dockerfile.Copy(),
			projpatterns.Readme.Copy(),
			projpatterns.GitIgnore.Copy(),
			projpatterns.Linter.Copy(),
			projpatterns.RscliMK.Copy(),
		}...,
	)
	return nil
}
func (a PrepareEnvironmentFoldersAction) NameInAction() string {
	return "Preparing environment folder"
}
