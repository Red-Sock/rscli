package config

import "github.com/Red-Sock/rscli/pkg/folder"

type EmptyConfig struct {
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

func (e *EmptyConfig) ExtractDataSources() (*folder.Folder, error) {
	return nil, nil
}

func (e *EmptyConfig) ExtractServerOptions() (*folder.Folder, error) {
	return nil, nil
}
