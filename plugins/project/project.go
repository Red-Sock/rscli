package project

import (
	"os"
	"strings"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

type Project struct {
	Name string
	Path string

	Cfg *config.Config

	ProjType Type

	Root folder.Folder
}

func (p *Project) GetShortName() string {
	name := p.Name

	if idx := strings.LastIndex(name, string(os.PathSeparator)); idx != -1 {
		name = name[idx+1:]
	}

	return name
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetFolder() *folder.Folder {
	return &p.Root
}

func (p *Project) GetConfig() *config.Config {
	return p.Cfg
}

func (p *Project) GetProjectPath() string {
	return p.Path
}

func (p *Project) GetType() Type {
	return p.ProjType
}
