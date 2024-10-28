package project_mock

import (
	"os"
	"path"
	"testing"

	"github.com/godverv/matreshka/environment"

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

func WithFileSystem(t *testing.T) Opt {
	return func(m *MockProject) {
		m.Path = path.Join(os.TempDir(), t.Name())
	}
}
