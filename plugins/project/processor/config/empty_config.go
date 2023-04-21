package config

import (
	"github.com/Red-Sock/rscli/pkg/folder"
)

type EmptyConfig struct {
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

func (e *EmptyConfig) GetServerOptions() ([]ServerOptions, error) {
	return nil, nil
}
