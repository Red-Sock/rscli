package project_mock

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
)

const testFolder = "test"

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
		m.Path = path.Join(testFolder, t.Name()[5:])
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

func WithSqlite(name string) Opt {
	return func(m *MockProject) {
		s := resources.NewSqlite(resources.Name(resources.SqliteResourceName + "_" + name))
		sq := s.(*resources.Sqlite)
		sq.Path = path.Join(sq.Path, name+".db")

		m.Cfg.DataSources = append(m.Cfg.DataSources, sq)
	}
}

func WithGrpcServer(port int) Opt {
	return func(m *MockProject) {
		m.Cfg.Servers[port] = &server.Server{
			Name: "MASTER",
			Port: strconv.Itoa(port),
			GRPC: map[string]*server.GRPC{
				"/": {
					Module:  m.Name,
					Gateway: "/api",
				},
			},
		}
	}
}
