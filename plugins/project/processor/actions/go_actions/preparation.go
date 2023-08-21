package go_actions

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type PrepareProjectStructureAction struct {
}

func (a PrepareProjectStructureAction) Do(p interfaces.Project) error {
	cmd := &folder.Folder{Name: patterns.CmdFolder}

	cmd.Inner = append(cmd.Inner, &folder.Folder{
		Name: p.GetName(),
		Inner: []*folder.Folder{
			{
				Name:    patterns.MainFileName,
				Content: patterns.MainFile,
			},
		},
	})

	fldr := p.GetFolder()
	fldr.Add(cmd)

	fldr.Add(&folder.Folder{Name: patterns.ConfigsFolder})

	fldr.Add(&folder.Folder{Name: patterns.InternalFolder})

	fldr.Add(&folder.Folder{
		Name: patterns.PkgFolder,
		Inner: []*folder.Folder{
			{Name: patterns.SwaggerFolder},
			{Name: patterns.ApiFolder},
		},
	})

	fldr.AddWithPath([]string{patterns.InternalFolder, patterns.UtilsFolder, patterns.CloserFolder}, &folder.Folder{
		Name:    patterns.CloserFile,
		Content: patterns.UtilsCloserFile,
	})

	return nil
}
func (a PrepareProjectStructureAction) NameInAction() string {
	return "Preparing project structure"
}

type PrepareExamplesFoldersAction struct{}

func (a PrepareExamplesFoldersAction) Do(p interfaces.Project) error {

	if p.GetFolder().GetByPath(patterns.ExamplesFolder, patterns.ExamplesHttpEnvFile) != nil {
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

	servers := p.GetConfig().GetServerOptions()

	for _, item := range servers {
		e.Dev[item.Name] = "0.0.0.0:" + strconv.FormatUint(uint64(item.Port), 10)
		e.DevDocker[item.Name] = "0.0.0.0:1" + strconv.FormatUint(uint64(item.Port), 10)
	}

	exampleFile, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		return errors.Wrap(err, "error marshalling example file")
	}

	p.GetFolder().Add(&folder.Folder{
		Name: patterns.ExamplesFolder,
		Inner: []*folder.Folder{
			{
				Name:    patterns.ExampleFileName,
				Content: patterns.ApiHTTP,
			},
			{
				Name:    patterns.ExamplesHttpEnvFile,
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
				Name:    patterns.DockerfileFileName,
				Content: patterns.Dockerfile,
			},
			{
				Name:    patterns.ReadMeFileName,
				Content: bytes.ReplaceAll(patterns.Readme, []byte("{{PROJECT_NAME}}"), []byte(p.GetName())),
			},
			{
				Name:    patterns.GitignoreFileName,
				Content: patterns.GitIgnore,
			},
			{
				Name:    patterns.GolangCIYamlFileName,
				Content: patterns.Linter,
			},
			{
				Name:    patterns.RsCliMkFileName,
				Content: patterns.RscliMK,
			},
		}...,
	)
	return nil
}
func (a PrepareEnvironmentFoldersAction) NameInAction() string {
	return "Preparing environment folder"
}
