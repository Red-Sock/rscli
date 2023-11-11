package go_actions

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p interfaces.Project) (err error) {
	configFolderPath := path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder)
	p.GetFolder().Add(&folder.Folder{
		Name: configFolderPath,
	})

	goConfigFolder := p.GetFolder().GetByPath(configFolderPath)

	err = a.generateGoKeysFile(p, goConfigFolder)
	if err != nil {
		return errors.Wrap(err, "error generating keys go-file")
	}

	a.generateGoConfigFiles(goConfigFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return errors.Wrap(err, "error generating config yaml-files")
	}

	return nil
}
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareGoConfigFolderAction) generateGoKeysFile(p interfaces.Project, goConfigFolder *folder.Folder) error {
	keys, err := matreshka.GenerateEnvironmentKeys(*p.GetConfig().AppConfig)
	if err != nil {
		return errors.Wrap(err, "error generating environment keys")
	}

	if len(keys) == 0 {
		goConfigFolder.GetByPath(projpatterns.ConfigKeysFileName).Delete()
		return nil
	}

	sb := rw.RW{}
	_, _ = sb.WriteString("package config\n\nconst (\n")
	for _, key := range keys {
		_ = sb.WriteByte('\t')
		_, _ = sb.WriteString(cases.SnakeToPascal(key.Name))
		_ = sb.WriteByte('=')
		_, _ = sb.WriteString(key.Name)
		_ = sb.WriteByte('\n')
	}
	_ = sb.WriteByte(')')

	cfgKeysFile := &folder.Folder{
		Name:    projpatterns.ConfigKeysFileName,
		Content: sb.Bytes(),
	}

	goConfigFolder.Add(cfgKeysFile)

	return nil
}

func (a PrepareGoConfigFolderAction) generateGoConfigFiles(goConfigFolder *folder.Folder) {
	goConfigFolder.Add(projpatterns.ConfigFile.Copy())
	goConfigFolder.Add(projpatterns.AutoloadConfigFile.Copy())
	goConfigFolder.Add(projpatterns.StaticConfigFile.Copy())
}

func (a PrepareGoConfigFolderAction) generateConfigYamlFile(p interfaces.Project) error {
	b, err := p.GetConfig().Marshal()
	if err != nil {
		return err
	}

	configFolder := p.GetFolder().GetByPath(projpatterns.ConfigsFolder)

	configFolder.Add(&folder.Folder{
		Name:    projpatterns.ConfigTemplate,
		Content: b,
	})
	configFolder.Add(&folder.Folder{
		Name:    projpatterns.DevConfigYamlFile,
		Content: b,
	})

	return nil
}
