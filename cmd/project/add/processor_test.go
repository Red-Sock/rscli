package project

import (
	"testing"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/server"
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

func Test_AddDependency(t *testing.T) {
	t.Parallel()

	type testCase struct {
		printlnCalls            []string
		cfg                     *config.RsCliConfig
		args                    []string
		expectedFiles           []*folder.Folder
		expectedMatreshkaServer *server.Server
	}

	testCases := map[string]testCase{
		"GRPC": {
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
			expectedMatreshkaServer: &server.Server{
				GRPC: make(map[string]*server.GRPC),
				FS:   make(map[string]*server.FS),
				HTTP: make(map[string]*server.HTTP),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name,
			func(t *testing.T) {
				tc := tc
				projectMock := project_mock.GetMockProject(t,
					project_mock.WithBasicConfig(t),
					project_mock.WithFileSystem(t),
					project_mock.WithGit(t),
				)
				defer func() {
					//require.NoError(t, os.RemoveAll(projectMock.Path))
				}()
				initActions := actions.InitProject(project.TypeGo)
				for _, action := range initActions {
					require.NoError(t, action.Do(projectMock))
				}

				ioMock := mocks.NewIOMock(t)
				setupPrintlnMock(t, ioMock, tc.printlnCalls...)

				apIoMock := mocks.NewIOMock(t)
				apIoMock.PrintlnMock.Set(func(_ ...string) {})

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

				command := Proc{
					Processor: processor.New(
						processor.WithIo(ioMock),
						processor.WithWd(projectMock.Project.GetProjectPath()),
						processor.WithConfig(tc.cfg),
					),
					ActionPerformer: apMock,
				}

				require.NoError(t, command.run(nil, tc.args))

				for _, expFile := range tc.expectedFiles {
					tests.AssertFolderInFs(t, projectMock.Path, expFile)
				}
			})
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
