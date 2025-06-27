package add

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/rerrors"

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
	cfg  *config.RsCliConfig
	args []string
}

// deleteGenerated - for debug purposes. Set to true when need to look at what script generates
const (
	deleteGenerated = false
	snapshotPath    = "./expected/"
)

func Test_AddDependency(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct{ prep func(t *testing.T) testCase }{
		"grpc": {
			prep: expectedGrpc,
		},
		"redis": {
			prep: expectedRedis,
		},
		"postgres": {
			prep: expectedPostgres,
		},
		"telegram": {
			prep: expectedTelegram,
		},
		"sqlite": {
			prep: expectedSqlite,
		},
		"env": {
			prep: expectedEnv,
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

				projectMock.Root.GetByPath(patterns.DockerfileFile).Delete()
				require.NoError(t, projectMock.Root.Build())

				ioMock := mocks.NewIOMock(t)
				setupPrintlnMock(t, ioMock, preparingMsg, startingMsg, endMsg)

				command := Proc{
					Processor: processor.New(
						processor.WithIo(ioMock),
						processor.WithWd(projectMock.Project.GetProjectPath()),
						processor.WithConfig(tc.cfg),
					),
					ActionPerformer: actions.NewActionPerformer(mocks.IoDevNul{}), //apMock(t),
				}
				require.NoError(t, command.run(nil, tc.args))

				expected, err := folder.Load(snapshotPath + name)
				require.NoError(t, err, "error loading expected folder")
				expected.Name = ""

				tests.AssertFolderInFs(t, projectMock.Path, expected)
			})
	}
}

func expectedGrpc(t *testing.T) testCase {

	return testCase{
		cfg: &config.RsCliConfig{
			Env: config.Project{
				PathToServerDefinition: "api",
			},
		},
		args: []string{dependencies.DependencyNameGrpc},
	}
}

func expectedRedis(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		args: []string{dependencies.DependencyNameRedis},
	}
}

func expectedPostgres(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		args: []string{dependencies.DependencyNamePostgres},
	}
}

func expectedTelegram(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		args: []string{dependencies.DependencyNameTelegram},
	}
}

func expectedSqlite(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		args: []string{dependencies.DependencyNameSqlite},
	}
}

func expectedEnv(t *testing.T) testCase {
	apIoMock := mocks.NewIOMock(t)
	apIoMock.PrintlnMock.Set(func(_ ...string) {})

	return testCase{
		args: []string{dependencies.DependencyEnvVariable},
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

func apMock(t *testing.T) actions.ActionPerformer {
	apm := mocks.NewActionPerformerMock(t)
	apm.TidyMock.Set(func(p project.IProject) (_ error) {
		shorActionList := []actions.Action{
			go_actions.PrepareConfigFolder{},
			go_actions.PrepareServer{},
			go_actions.BuildProjectAction{},
		}

		for _, a := range shorActionList {
			err := a.Do(p)
			if err != nil {
				return rerrors.Wrap(err, "error performing action: ", a.NameInAction())
			}
		}

		return nil
	})

	return apm
}
