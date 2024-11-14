package project_mock

import (
	"os"
	"testing"

	"github.com/godverv/matreshka/environment"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
)

const testFolder = "test_"

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
		m.Path = testFolder + t.Name()[5:]
		m.Root.Name = m.Path
		require.NoError(t, os.MkdirAll(m.Path, 0777))
	}
}

func WithBasicConfig(t *testing.T) Opt {
	return func(m *MockProject) {
		require.NoError(t, m.Cfg.Unmarshal(BasicConfig()))
	}
}

func WithGit(t *testing.T) Opt {
	return func(m *MockProject) {
		require.NotEmpty(t, m.Path, "to enable git in mock project WithFileSystem is required")
		require.NoError(t, git.Init(m.Project.GetProjectPath()))
	}
}
