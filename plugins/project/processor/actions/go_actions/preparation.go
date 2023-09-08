package go_actions

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	patterns2 "github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type PrepareProjectStructureAction struct {
}

func (a PrepareProjectStructureAction) Do(p interfaces.Project) error {
	cmd := &folder.Folder{Name: patterns2.CmdFolder}

	cmd.Inner = append(cmd.Inner, &folder.Folder{
		Name: p.GetShortName(),
		Inner: []*folder.Folder{
			{
				Name:    patterns2.MainFileName,
				Content: patterns2.MainFile,
			},
		},
	})

	fldr := p.GetFolder()
	fldr.Add(cmd)

	fldr.Add(&folder.Folder{Name: patterns2.ConfigsFolder})

	fldr.Add(&folder.Folder{Name: patterns2.InternalFolder})

	fldr.Add(&folder.Folder{
		Name: patterns2.PkgFolder,
		Inner: []*folder.Folder{
			{Name: patterns2.SwaggerFolder},
			{Name: patterns2.ApiFolder},
		},
	})

	fldr.AddWithPath([]string{patterns2.InternalFolder, patterns2.UtilsFolder, patterns2.CloserFolder}, &folder.Folder{
		Name:    patterns2.CloserFile,
		Content: patterns2.UtilsCloserFile,
	})

	return nil
}
func (a PrepareProjectStructureAction) NameInAction() string {
	return "Preparing project structure"
}

type PrepareExamplesFoldersAction struct{}

func (a PrepareExamplesFoldersAction) Do(p interfaces.Project) error {

	if p.GetFolder().GetByPath(patterns2.ExamplesFolder, patterns2.ExamplesHttpEnvFile) != nil {
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
		Name: patterns2.ExamplesFolder,
		Inner: []*folder.Folder{
			{
				Name:    patterns2.ExampleFileName,
				Content: patterns2.ApiHTTP,
			},
			{
				Name:    patterns2.ExamplesHttpEnvFile,
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
				Name:    patterns2.DockerfileFileName,
				Content: patterns2.Dockerfile,
			},
			{
				Name:    patterns2.ReadMeFileName,
				Content: bytes.ReplaceAll(patterns2.Readme, []byte("{{PROJECT_NAME}}"), []byte(p.GetName())),
			},
			{
				Name:    patterns2.GitignoreFileName,
				Content: patterns2.GitIgnore,
			},
			{
				Name:    patterns2.GolangCIYamlFileName,
				Content: patterns2.Linter,
			},
			{
				Name:    patterns2.RsCliMkFileName,
				Content: patterns2.RscliMK,
			},
		}...,
	)
	return nil
}
func (a PrepareEnvironmentFoldersAction) NameInAction() string {
	return "Preparing environment folder"
}
