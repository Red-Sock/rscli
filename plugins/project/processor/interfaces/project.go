package interfaces

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
)

type Project interface {
	GetName() string
	GetConfig() ProjectConfig
	GetProjectPath() string

	GetFolder() *folder.Folder
	GetVersion() Version
	SetVersion(Version)
}

type ProjectConfig interface {
	GetPath() string
	SetPath(pth string)

	GetTemplate() ([]byte, error)

	GenerateGoConfigKeys(projName string) ([]byte, error)

	ExtractName() (string, error)
	GetDataSourceFolders() (*folder.Folder, error)
	GetServerFolders() ([]*folder.Folder, error)
	GetServerOptions() ([]config.ServerOptions, error)
}
