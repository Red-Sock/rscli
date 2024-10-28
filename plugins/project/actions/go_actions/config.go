package go_actions

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/config_generators"
)

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p project.IProject) (err error) {
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
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareGoConfigFolderAction) generateConfigYamlFile(p project.IProject) error {
	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)

	newConfig := p.GetConfig()
	// Dev config
	{
		err := appendToConfig(newConfig.AppConfig, configFolder, patterns.ConfigDevYamlFile)
		if err != nil {
			return errors.Wrap(err, "error appending changes to dev config")
		}
	}

	obfuscateConfig(p.GetConfig())

	// Template
	{
		err := appendToConfig(newConfig.AppConfig, configFolder, patterns.ConfigTemplateYaml)
		if err != nil {
			return errors.Wrap(err, "error appending changes to dev config")
		}
	}

	// Master config
	{
		err := appendToConfig(newConfig.AppConfig, configFolder, patterns.ConfigMasterYamlFile)
		if err != nil {
			return errors.Wrap(err, "error appending changes to dev config")
		}
	}

	return nil
}

func obfuscateConfig(cfg *config.Config) {
	for i := range cfg.DataSources {
		cfg.DataSources[i] = cfg.DataSources[i].Obfuscate()
	}
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
