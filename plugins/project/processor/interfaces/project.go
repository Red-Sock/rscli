package interfaces

import "github.com/Red-Sock/rscli/pkg/folder"

type Project interface {
	GetName() string
	GetConfig() Config
	GetProjectPath() string

	GetFolder() *folder.Folder
}

type Config interface {
	GetPath() string
	SetPath(pth string)

	GenerateGoConfigKeys(projName string) ([]byte, error)

	ExtractName() (string, error)
	ExtractDataSources() (*folder.Folder, error)
	ExtractServerOptions() ([]*folder.Folder, error)
}
