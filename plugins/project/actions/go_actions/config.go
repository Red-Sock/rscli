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
	cfgFolder, err := config_generators.GenerateConfigFolder(p.GetConfig())
	if err != nil {
		return errors.Wrap(err, "error generating config folder")
	}
	cfgFolder.Name = path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder)

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

func obfuscateConfig(cfg *config.Config) {
	for i := range cfg.DataSources {
		cfg.DataSources[i] = cfg.DataSources[i].Obfuscate()
	}
}
