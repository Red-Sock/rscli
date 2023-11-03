package go_actions

import (
	"path"

	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p interfaces.Project) error {
	configFolderPath := path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder)
	p.GetFolder().Add(&folder.Folder{
		Name: configFolderPath,
	})

	configFolder := p.GetFolder().GetByPath(configFolderPath)
	configFolder.Add(projpatterns.ConfigFile.Copy())

	keys, err := matreshka.GenerateGoConfigKeys(p.GetShortName(), p.GetConfig().AppConfig)
	if err != nil {
		return err
	}

	cfgKeysFile := projpatterns.ConfigKeysFile.Copy()

	cfgKeysFile.Content = append(cfgKeysFile.Content, []byte("const (\n")...)
	cfgKeysFile.Content = append(cfgKeysFile.Content, keys...)
	cfgKeysFile.Content = append(cfgKeysFile.Content, ')')

	configFolder.Add(cfgKeysFile)

	return nil
}
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

type BuildProjectAction struct{}

func (a BuildProjectAction) Do(p interfaces.Project) error {
	ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build()
	if err != nil {
		return err
	}
	return nil
}
func (a BuildProjectAction) NameInAction() string {
	return "Building project"
}
