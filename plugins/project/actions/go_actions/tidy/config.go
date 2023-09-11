package tidy

import (
	"path"

	"github.com/Red-Sock/trace-errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

var ErrNoMakeFile = errors.New("no rscli.mk makefile found")

// TODO
func Config(p interfaces.Project) error {
	config := p.GetConfig()

	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	p.GetFolder().ForceAddWithPath([]string{patterns.ConfigsFolder}, &folder.Folder{
		Name:    patterns.ConfigTemplate,
		Content: b,
	})

	appInfo := config.GetProjInfo()

	if appInfo.Name != "" {
		modName := p.GetName()

		if modName != appInfo.Name {
			appInfo.Name = modName
			//todo change to path to project + path to config
			err = config.BuildTo(path.Join(p.GetProjectPath(), patterns.ConfigsFolder, patterns.ConfigYamlFile))
			if err != nil {
				return errors.Wrap(err, "error during rebuilding")
			}
		}
	}

	return nil
}
