package go_actions

import (
	"encoding/json"
	"path"
	"strconv"
	"strings"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	patterns "github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type PrepareProjectStructureAction struct {
}

func (a PrepareProjectStructureAction) Do(p interfaces.Project) error {
	rootF := p.GetFolder()

	{
		cmd := &folder.Folder{Name: patterns.CmdFolder}
		mainFilePath := path.Join(strings.ToLower(p.GetShortName()), patterns.MainFile.Name)

		cmd.Add(patterns.MainFile.CopyWithNewName(mainFilePath))

		rootF.Add(cmd)
	}

	{
		rootF.Add(&folder.Folder{Name: patterns.ConfigsFolder})
		rootF.Add(&folder.Folder{Name: patterns.InternalFolder})
	}

	{
		rootF.Add(&folder.Folder{
			Name:  patterns.PkgFolder,
			Inner: []*folder.Folder{},
		})

		closerFilePath := path.Join(patterns.InternalFolder, patterns.UtilsFolder, patterns.CloserFolder, patterns.UtilsCloserFile.Name)
		rootF.Add(patterns.UtilsCloserFile.CopyWithNewName(closerFilePath))
	}
	return nil
}
func (a PrepareProjectStructureAction) NameInAction() string {
	return "Preparing project structure"
}

// TODO RSI-245 - мб переделать на общие какие-то запросы
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

	for _, item := range p.GetConfig().Servers {
		portStr := strconv.FormatUint(uint64(item.GetPort()), 10)
		e.Dev[item.GetName()] = "0.0.0.0:" + portStr
		e.DevDocker[item.GetName()] = "0.0.0.0:1" + portStr
	}

	exampleFile, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		return errors.Wrap(err, "error marshalling example file")
	}

	p.GetFolder().Add(&folder.Folder{
		Name: patterns.ExamplesFolder,
		Inner: []*folder.Folder{
			patterns.ApiHTTP.Copy(),
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
		patterns.Dockerfile.Copy(),
		patterns.Readme.Copy(),
		patterns.GitIgnore.Copy(),
		patterns.Linter.Copy(),
	)

	return nil
}
func (a PrepareEnvironmentFoldersAction) NameInAction() string {
	return "Preparing environment folder"
}
