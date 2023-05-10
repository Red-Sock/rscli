package interfaces

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/config/pkg/configstructs"
)

type Project interface {
	GetName() string
	GetProjectModName() string
	GetConfig() ProjectConfig
	GetProjectPath() string

	GetFolder() *folder.Folder
	GetVersion() Version
	SetVersion(Version)
}

type ProjectConfig interface {
	Rebuild(p Project) error

	GetPath() string
	SetPath(pth string)
	GetProjInfo() (*configstructs.AppInfo, error)

	GetTemplate() ([]byte, error)

	GenerateGoConfigKeys(projName string) ([]byte, error)

	ExtractName() (string, error)
	GetDataSourceFolders() (*folder.Folder, error)
	GetServerFolders() ([]*folder.Folder, error)
	GetServerOptions() ([]configstructs.ServerOptions, error)
	GetDataSourceOptions() (out []configstructs.ConnectionOptions, err error)
}
