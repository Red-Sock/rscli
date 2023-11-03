package tidy

import (
	"os"
	"path"

	"github.com/Red-Sock/trace-errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

var ErrNoMakeFile = errors.New("no rscli.mk makefile found")

// TODO
func Config(p interfaces.Project) error {
	config := p.GetConfig()

	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	p.GetFolder().Add(&folder.Folder{
		Name:    path.Join(projpatterns.ConfigsFolder, projpatterns.ConfigTemplate),
		Content: b,
	})

	p.GetFolder().Add(&folder.Folder{
		Name:    path.Join(projpatterns.ConfigsFolder, projpatterns.DevConfigYamlFile),
		Content: b,
	})

	if config.AppInfo.Name != "" {
		modName := p.GetName()

		if modName != config.AppInfo.Name {
			config.AppInfo.Name = modName
			cfgPath := path.Join(p.GetProjectPath(), projpatterns.ConfigsFolder, projpatterns.ConfigYamlFile)
			b, err := p.GetConfig().Marshal()
			if err != nil {
				return errors.Wrap(err, "error marshalling config")
			}

			err = os.WriteFile(cfgPath, b, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "error writing config file")
			}
		}
	}

	return nil
}
