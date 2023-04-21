package tidy

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

func Config(p interfaces.Project) error {
	config := p.GetConfig()
	b, err := config.GetTemplate()
	if err != nil {
		return err
	}

	tmplt := p.GetFolder().GetByPath(patterns.ConfigsFolder, patterns.ConfigTemplate)
	if tmplt == nil {
		p.GetFolder().AddWithPath([]string{patterns.ConfigsFolder}, &folder.Folder{
			Name:    patterns.ConfigTemplate,
			Content: b,
		})
	} else {
		tmplt.Content = b
	}

	return nil
}
