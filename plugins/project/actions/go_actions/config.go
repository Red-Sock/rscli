package go_actions

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/config_generators"
)

type PrepareConfigFolder struct{}

func (a PrepareConfigFolder) Do(p project.IProject) (err error) {
	cfgFolder, err := config_generators.GenerateConfigFolder(p.GetConfig())
	if err != nil {
		return errors.Wrap(err, "error generating config folder")
	}

	cfgFolder.Name = path.Join(patterns.InternalFolder, patterns.ConfigsFolder)

	p.GetFolder().Add(cfgFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return errors.Wrap(err, "error generating config yaml-files")
	}

	return nil
}
func (a PrepareConfigFolder) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareConfigFolder) generateConfigYamlFile(p project.IProject) error {
	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)

	newConfig := p.GetConfig()

	for _, cfgName := range []string{
		patterns.ConfigDevYamlFile,
		patterns.ConfigTemplateYaml,
		patterns.ConfigMasterYamlFile,
	} {
		err := appendToConfig(newConfig.AppConfig, configFolder, cfgName)
		if err != nil {
			return errors.Wrap(err, "error appending changes to dev config")
		}
	}

	return nil
}

func appendToConfig(newConfig matreshka.AppConfig, configFolder *folder.Folder, path string) (err error) {
	currentConfig := matreshka.NewEmptyConfig()

	configFile := configFolder.GetByPath(path)
	if configFile == nil {
		configFile = &folder.Folder{
			Name: path,
		}
		configFolder.Add(configFile)
	}

	if len(configFile.Content) != 0 {
		err = currentConfig.Unmarshal(configFile.Content)
		if err != nil {
			return errors.Wrap(err, "error reading dev config file")
		}
	}

	currentConfig = matreshka.MergeConfigs(currentConfig, newConfig)

	configFile.Content, err = currentConfig.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling dev config to yaml")
	}

	return nil
}
