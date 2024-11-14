package add

import (
	"os"
	"testing"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests"
	"github.com/Red-Sock/rscli/tests/mocks"
	"github.com/Red-Sock/rscli/tests/project_mock"
)

type testCase struct {
	actionPerformer actions.ActionPerformer
	printlnCalls    []string
	cfg             *config.RsCliConfig
	args            []string
	expectedFiles   []*folder.Folder
}

// deleteGenerated - for debug purposes. Set to true when need to look at what script generates
const deleteGenerated = false

func Test_AddDependency(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct{ prep func(t *testing.T) testCase }{
		"GRPC": {
			prep: expectedGrpc,
		},
		"REDIS": {
			prep: expectedRedis,
		},
		"POSTGRES": {
			prep: expectedPostgres,
		},
		"TELEGRAM": {
			prep: expectedTelegram,
		},
		"SQLITE": {
			prep: expectedSqlite,
		},
	}

	for name, tc := range testCases {
		t.Run(name,
			func(t *testing.T) {
				t.Parallel()
				tc := tc.prep(t)
				projectMock := project_mock.GetMockProject(t,
					project_mock.WithBasicConfig(t),
					project_mock.WithFileSystem(t),
					project_mock.WithGit(t),
				)
				defer func() {
					if deleteGenerated {
						require.NoError(t, os.RemoveAll(projectMock.Path))
					}
				}()
				initActions := actions.InitProject(project.TypeGo)
				for _, action := range initActions {
					require.NoError(t, action.Do(projectMock))
				}

				ioMock := mocks.NewIOMock(t)
				setupPrintlnMock(t, ioMock, tc.printlnCalls...)

				command := Proc{
					Processor: processor.New(
						processor.WithIo(ioMock),
						processor.WithWd(projectMock.Project.GetProjectPath()),
						processor.WithConfig(tc.cfg),
					),
					ActionPerformer: tc.actionPerformer,
				}

				require.NoError(t, command.run(nil, tc.args))

				for _, expFile := range tc.expectedFiles {
					tests.AssertFolderInFs(t, projectMock.Path, expFile)
				}
			})
	}
}

func expectedGrpc(t *testing.T) testCase {
	apMock := mocks.NewActionPerformerMock(t)
	apMock.TidyMock.Set(func(p project.IProject) (_ error) {
		shorActionList := []actions.Action{
			go_actions.PrepareGoConfigFolderAction{},
			go_actions.PrepareServerAction{},
			go_actions.BuildProjectAction{},
		}

		for _, a := range shorActionList {
			err := a.Do(p)
			if err != nil {
				return errors.Wrap(err, "error performing action: ", a.NameInAction())
			}
		}

		return nil
	})

	return testCase{
		actionPerformer: apMock,
		printlnCalls: []string{
			preparingMsg,
			startingMsg,
			endMsg,
		},
		cfg: &config.RsCliConfig{
			Env: config.Project{
				PathToServerDefinition: "api",
			},
		},
		args: []string{dependencies.DependencyNameGrpc},
		expectedFiles: []*folder.Folder{
			{
				Name:    "api/grpc/GRPC_api.proto",
				Content: grpcExpectedProtoFile,
			},
			{
				Name: "config",
				Inner: []*folder.Folder{
					{
						Name:    "config.yaml",
						Content: grpcMatreshkaConfigExpected,
					},
					{
						Name:    "config_template.yaml",
						Content: grpcMatreshkaConfigExpected,
					},
					{
						Name:    "dev.yaml",
						Content: grpcMatreshkaConfigExpected,
					},
				},
			},
			{
				Name: "internal/transport",
				Inner: []*folder.Folder{
					{
						Name:    "grpc.go",
						Content: patterns.GrpcServerManager.Content,
					},
					{
						Name:    "http.go",
						Content: patterns.HttpServerManager.Content,
					},
					{
						Name:    "manager.go",
						Content: patterns.ServerManager.Content,
					},
				},
			},
		},
	}
}

func expectedRedis(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		actionPerformer: actions.NewActionPerformer(apIoMock),
		printlnCalls: []string{
			preparingMsg,
			startingMsg,
			endMsg,
		},
		args: []string{dependencies.DependencyNameRedis},
		expectedFiles: []*folder.Folder{
			{
				Name: "config",
				Inner: []*folder.Folder{
					{
						Name:    "config.yaml",
						Content: expectedRedisConfig,
					},
					{
						Name:    "config_template.yaml",
						Content: expectedRedisConfig,
					},
					{
						Name:    "dev.yaml",
						Content: expectedRedisConfig,
					},
				},
			},
			{
				Name:    "internal/clients/redis/conn.go",
				Content: patterns.RedisConnFile.Content,
			},
			{
				Name:    "internal/config/data_sources.go",
				Content: expectedRedisDataSourceConfig,
			},
			{
				Name:    "internal/app/data_sources.go",
				Content: expectedRedisDataSourceApp,
			},
		},
	}
}

func expectedPostgres(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		actionPerformer: actions.NewActionPerformer(apIoMock),
		printlnCalls: []string{
			preparingMsg,
			startingMsg,
			endMsg,
		},
		args: []string{dependencies.DependencyNamePostgres},
		expectedFiles: []*folder.Folder{
			{
				Name: "config",
				Inner: []*folder.Folder{
					{
						Name:    "config.yaml",
						Content: expectedPostgresConfig,
					},
				},
			},
			{
				Name:    "internal/app/data_sources.go",
				Content: expectedPostgresDataSourceApp,
			},
			{
				Name:    "internal/config/data_sources.go",
				Content: expectedPostgresDataSourceConfig,
			},
			{
				Name:    "internal/clients/sqldb/conn.go",
				Content: patterns.SqlConnFile.Content,
			},
			{
				Name:    "internal/clients/sqldb/postgres.go",
				Content: patterns.PostgresConnFile.Content,
			},
		},
	}
}

func expectedTelegram(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		actionPerformer: actions.NewActionPerformer(apIoMock),
		printlnCalls: []string{
			preparingMsg,
			startingMsg,
			endMsg,
		},
		args: []string{dependencies.DependencyNameTelegram},
		expectedFiles: []*folder.Folder{
			{
				Name:    "config/config.yaml",
				Content: expectedTelegramConfig,
			},
			{
				Name:    "internal/app/data_sources.go",
				Content: expectedTelegramDataSourcesApp,
			},
			{
				Name:    "internal/config/data_sources.go",
				Content: expectedTelegramDataSourcesConfig,
			},
			{
				Name:    "internal/clients/telegram/conn.go",
				Content: patterns.TgConnFile.Content,
			},
			{
				Name:    "internal/transport/telegram/listener.go",
				Content: expectedTelegramServer,
			},
			{
				Name:    "internal/transport/telegram/version/handler.go",
				Content: expectedTelegramServerHandlerExample,
			},
		},
	}
}

func expectedSqlite(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		actionPerformer: actions.NewActionPerformer(apIoMock),
		printlnCalls: []string{
			preparingMsg,
			startingMsg,
			endMsg,
		},
		args: []string{dependencies.DependencyNameSqlite},
		expectedFiles: []*folder.Folder{
			{
				Name:    "config/config.yaml",
				Content: expectedSqliteConfig,
			},
			{
				Name:    "internal/app/data_sources.go",
				Content: expectedSqliteDataSourcesApp,
			},
			{
				Name:    "internal/config/data_sources.go",
				Content: expectedSqliteDataSourcesConfig,
			},
			{
				Name:    "internal/clients/sqldb/conn.go",
				Content: patterns.SqlConnFile.Content,
			},
			{
				Name:    "internal/clients/sqldb/sqlite.go",
				Content: patterns.SqliteConnFile.Content,
			},
		},
	}
}

func setupPrintlnMock(t *testing.T, ioMock *mocks.IOMock, printlnCalls ...string) {
	prinlnIdx := 0
	ioMock.PrintlnMock.Set(func(in ...string) {
		if prinlnIdx >= len(printlnCalls) {
			return
		}
		for _, text := range in {
			require.Equal(t, printlnCalls[prinlnIdx], text)
			prinlnIdx++
		}
	})
}
