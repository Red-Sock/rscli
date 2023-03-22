package interfaces

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
)

type Project interface {
	GetName() string
	GetConfig() Config
	GetProjectPath() string

	GetFolder() *folder.Folder
	GetVersion() Version
	SetVersion(Version)
}

type Config interface {
	GetPath() string
	SetPath(pth string)

	GenerateGoConfigKeys(projName string) ([]byte, error)

	ExtractName() (string, error)
	GetDataSourceFolders() (*folder.Folder, error)
	GetServerFolders() ([]*folder.Folder, error)
	GetServerOptions() ([]config.ServerOptions, error)
}
