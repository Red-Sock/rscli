package go_actions

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators/config_generators"
)

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p project.Project) (err error) {
	configFolder := &folder.Folder{
		Name: path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder),
	}

	configFolder.Add(projpatterns.AutoloadConfigFile.Copy())

	configStructFolders, err := a.generateConfigStructsFiles(p.GetConfig())
	if err != nil {
		return errors.Wrap(err, "error generating config structs files")
	}
	configFolder.Add(configStructFolders...)

	p.GetFolder().Add(configFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return errors.Wrap(err, "error generating config yaml-files")
	}

	return nil
}
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareGoConfigFolderAction) generateConfigYamlFile(p project.Project) error {
	configFolder := p.GetFolder().GetByPath(projpatterns.ConfigsFolder)

	plainConfig, err := p.GetConfig().Marshal()
	if err != nil {
		return errors.Wrap(err)
	}
	configFolder.Add(
		&folder.Folder{
			Name:    projpatterns.DevConfigYamlFile,
			Content: plainConfig,
		})

	obfuscateConfig(p.GetConfig())

	protectedCfg, err := p.GetConfig().Marshal()
	if err != nil {
		return errors.Wrap(err)
	}

	prodConfig := configFolder.GetByPath(projpatterns.ConfigYamlFile)
	if prodConfig == nil {
		configFolder.Add(&folder.Folder{
			Name:    projpatterns.ConfigYamlFile,
			Content: protectedCfg,
		})
	}

	configFolder.Add(
		&folder.Folder{
			Name:    projpatterns.ConfigTemplate,
			Content: protectedCfg,
		})

	return nil
}

func (a PrepareGoConfigFolderAction) generateConfigStructsFiles(cfg *config.Config) ([]*folder.Folder, error) {
	out := make([]*folder.Folder, 0, 3)
	// Environment config
	{
		confStruct, err := config_generators.GenerateEnvironmentConfigStruct(cfg.Environment)
		if err != nil {
			return nil, errors.Wrap(err, "error generating environment struct file")
		}
		out = append(out,
			&folder.Folder{
				Name:    projpatterns.ConfigEnvironmentFileName,
				Content: confStruct,
			})
	}
	// Data sources
	{
		confStruct, err := config_generators.GenerateDataSourcesConfigStruct(cfg.DataSources)
		if err != nil {
			return nil, errors.Wrap(err, "error generating data sources struct file")
		}
		out = append(out,
			&folder.Folder{
				Name:    projpatterns.ConfigDataSourcesFileName,
				Content: confStruct,
			})
	}

	return out, nil
}

func obfuscateConfig(cfg *config.Config) {
	for i := range cfg.DataSources {
		cfg.DataSources[i] = cfg.DataSources[i].Obfuscate()
	}
}
