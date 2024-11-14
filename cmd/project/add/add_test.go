package project

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/tests/mocks"
	"github.com/Red-Sock/rscli/tests/project_mock"
)

func Test_AddDependency(t *testing.T) {
	t.Parallel()

	mp := project_mock.GetMockProject(t,
		project_mock.WithBasicConfig(t),
		project_mock.WithFileSystem(t),
		project_mock.WithGit(t),
	)
	defer func() {
		//require.NoError(t, os.RemoveAll(mp.Path))
	}()
	require.NoError(t, go_actions.BuildProjectAction{}.Do(mp))

	io := mocks.NewIOMock(t)

	{
		printlnCalls := []string{
			preparingMsg,
			startingMsg,
			endMsg,
		}
		prinlnIdx := 0
		io.PrintlnMock.Set(func(in ...string) {
			if prinlnIdx >= len(printlnCalls) {
				require.Fail(t, "unknown println call", in)
			}
			for _, text := range in {
				require.Equal(t, printlnCalls[prinlnIdx], text)
				prinlnIdx++
			}
		})
	}

	cfg := &config.RsCliConfig{}

	basicProc := processor.New(
		processor.WithIo(io),
		processor.WithWd(mp.Project.GetProjectPath()),
		processor.WithConfig(cfg),
	)

	cmd := NewCommand(basicProc)
	args := []string{"grpc"}
	cmd.SetArgs(args)

	require.NoError(t, cmd.Execute())
}
