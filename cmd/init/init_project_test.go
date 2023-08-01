package init

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/pkg/colors"
	"github.com/Red-Sock/rscli/tests/mocks"
)

func Test_InitProject(t *testing.T) {
	const hintMessage = `
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "github.com/RedSock" and final result is "github.com/RedSock/rscli"
>`

	t.Run("OK_SINGLE_NAME", func(t *testing.T) {
		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		ioMock.GetInputMock.Expect().Return("rscli", nil)

		ioMock.PrintlnColoredMock.Expect(colors.ColorCyan, `Wonderful!!! "github.com/RedSock/rscli" it is!`)

		p := projectConstructor{
			cfg: cfg,
			io:  ioMock,
		}

		err := p.initProject(nil, nil)
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})
	t.Run("OK_NAME_WITH_SHORT_URL", func(t *testing.T) {
		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		ioMock.GetInputMock.Expect().Return("gitlab.com/RedSock/rscli", nil)

		ioMock.PrintlnColoredMock.Expect(colors.ColorCyan, `Wonderful!!! "gitlab.com/RedSock/rscli" it is!`)

		p := projectConstructor{
			cfg: cfg,
			io:  ioMock,
		}

		err := p.initProject(nil, nil)
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})
	t.Run("OK_NAME_WITH_HTTP_URL", func(t *testing.T) {
		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		ioMock.GetInputMock.Expect().Return("https://gitlab.com/RedSock/rscli", nil)

		ioMock.PrintlnColoredMock.Expect(colors.ColorCyan, `Wonderful!!! "gitlab.com/RedSock/rscli" it is!`)

		p := projectConstructor{
			cfg: cfg,
			io:  ioMock,
		}

		err := p.initProject(nil, nil)
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

	t.Run("ERROR_GET_PROJECT_NAME", func(t *testing.T) {
		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		errWant := errors.New("input error")
		ioMock.GetInputMock.Expect().Return("", errWant)

		p := projectConstructor{
			cfg: cfg,
			io:  ioMock,
		}

		errGot := p.initProject(nil, nil)
		require.ErrorIs(t, errGot, errWant)

	})
	t.Run("ERROR_INVALID_NAME", func(t *testing.T) {
		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		ioMock.GetInputMock.Expect().Return("rscli$1", nil)

		ioMock.PrintlnColoredMock.Expect(colors.ColorCyan, `Wonderful!!! "gitlab.com/RedSock/rscli" it is!`)

		p := projectConstructor{
			cfg: cfg,
			io:  ioMock,
		}

		err := p.initProject(nil, nil)
		require.Contains(t, err.Error(), "name contains \"$\" symbol")
		ioMock.MinimockPrintlnInspect()
		ioMock.MinimockGetInputInspect()
	})
}
