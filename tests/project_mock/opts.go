package project_mock

import (
	"github.com/godverv/matreshka/environment"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

func WithFile(filePath string, file []byte) Opt {
	return func(m *MockProject) {
		m.Root.Add(&folder.Folder{
			Name:    filePath,
			Content: file,
		})
	}
}

func WithEnvironmentVariables(vars ...*environment.Variable) Opt {
	return func(m *MockProject) {
		m.Cfg.Environment = append(m.Cfg.Environment, vars...)
	}
}

func WithDataSources(ds ...resources.Resource) Opt {
	return func(m *MockProject) {
		m.Cfg.DataSources = append(m.Cfg.DataSources, ds...)
	}
}
