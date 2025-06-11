package dockerfile_generator

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/tests/project_mock"
)

var (
	//go:embed test_snapshots/basic.Dockerfile
	basicDockerfile string
	//go:embed test_snapshots/with_sqlite.Dockerfile
	sqliteDockerfile string
	//go:embed test_snapshots/with_sqlite_multiple.Dockerfile
	multipleSqliteDockerfile string

	//go:embed test_snapshots/with_server.Dockerfile
	serverDockerfile string
	//go:embed test_snapshots/with_multiple_servers.Dockerfile
	multipleServersDockerfile string
	//go:embed test_snapshots/with_multiple_servers_and_sqlite.Dockerfile
	multipleServersAndSqliteDockerfile string
)

func Test_GenerateDockerfile(t *testing.T) {

	type testCase struct {
		genProj  func() project.IProject
		expected string
	}

	testCases := map[string]testCase{
		"basic": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t)
			},
			expected: basicDockerfile,
		},

		"with_sqlite": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithSqlite("test"))
			},
			expected: sqliteDockerfile,
		},
		"with_multiple_sqlite": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithSqlite("test"),
					project_mock.WithSqlite("test2"),
				)
			},
			expected: multipleSqliteDockerfile,
		},

		"with_server": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithGrpcServer(50051))
			},
			expected: serverDockerfile,
		},
		"with_multiple_servers": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithGrpcServer(50051),
					project_mock.WithGrpcServer(50052))
			},
			expected: multipleServersDockerfile,
		},

		"with_multiple_servers_and_sqlite": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithGrpcServer(50051),
					project_mock.WithGrpcServer(50052),
					project_mock.WithSqlite("test"),
					project_mock.WithSqlite("test2"),
				)
			},
			expected: multipleServersAndSqliteDockerfile,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			proj := tc.genProj()

			file, err := GenerateDockerfile(proj)
			require.NoError(t, err)
			require.Equal(t, string(file), tc.expected)
		})
	}
}
