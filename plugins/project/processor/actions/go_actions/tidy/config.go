package tidy

import (
	"path"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	patterns2 "github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var ErrNoMakeFile = errors.New("no rscli.mk makefile found")

// TODO
func Config(p interfaces.Project) error {
	config := p.GetConfig()

	b, err := config.GetTemplate()
	if err != nil {
		return err
	}

	p.GetFolder().ForceAddWithPath([]string{patterns2.ConfigsFolder}, &folder.Folder{
		Name:    patterns2.ConfigTemplate,
		Content: b,
	})

	appInfo := config.GetProjInfo()

	if appInfo.Name != "" {
		modName := p.GetName()

		if modName != appInfo.Name {
			appInfo.Name = modName
			//todo change to path to project + path to config
			err = config.BuildTo(path.Join(p.GetProjectPath(), patterns2.ConfigsFolder, patterns2.ConfigYamlFile))
			if err != nil {
				return errors.Wrap(err, "error during rebuilding")
			}
		}
	}

	return nil
}
