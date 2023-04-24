package tidy

import (
	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var ErrNoMakeFile = errors.New("no rscli.mk makefile found")

func Config(p interfaces.Project) error {
	config := p.GetConfig()
	b, err := config.GetTemplate()
	if err != nil {
		return err
	}

	p.GetFolder().ForceAddWithPath([]string{patterns.ConfigsFolder}, &folder.Folder{
		Name:    patterns.ConfigTemplate,
		Content: b,
	})

	return nil
}
