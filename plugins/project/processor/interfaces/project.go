package interfaces

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/config/pkg/configstructs"
)

type Project interface {
	GetName() string
	GetShortName() string

	GetConfig() ProjectConfig

	GetFolder() *folder.Folder
	GetProjectPath() string

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
