package config

import (
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/config/pkg/configstructs"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

type EmptyConfig struct {
}

func (e *EmptyConfig) Rebuild(_ interfaces.Project) error {
	return nil
}

func (e *EmptyConfig) GetProjInfo() (*configstructs.AppInfo, error) {
	return nil, nil
}

func (e *EmptyConfig) GetDataSourceOptions() (out []configstructs.ConnectionOptions, err error) {
	return nil, nil
}

func (e *EmptyConfig) GetTemplate() ([]byte, error) {
	return nil, nil
}

func NewEmptyProjectConfig() *EmptyConfig {
	return &EmptyConfig{}
}

func (e *EmptyConfig) GetPath() string {
	return ""
}

func (e *EmptyConfig) SetPath(pth string) {
}

func (e *EmptyConfig) GenerateGoConfigKeys(projName string) ([]byte, error) {
	return nil, nil
}

func (e *EmptyConfig) GetDataSourceFolders() (*folder.Folder, error) {
	return nil, nil
}

func (e *EmptyConfig) GetServerFolders() ([]*folder.Folder, error) {
	return nil, nil
}

func (e *EmptyConfig) ExtractName() (string, error) {
	return "", nil
}

func (e *EmptyConfig) GetServerOptions() ([]configstructs.ServerOptions, error) {
	return nil, nil
}
