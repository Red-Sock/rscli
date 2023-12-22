package project

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/godverv/matreshka/resources"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
	"github.com/Red-Sock/rscli/tests/mocks"
)

type projectAddValidator interface {
	getName() string
	validate(t *testing.T, pth string)
}

const (
	stdConfigPathToClient = "internal/clients"
	stdServerPath         = "internal/transport"
	stdConfigPath         = "config/config.yaml"
	projectName           = "rscli"
)

func Test_ProjectAdd(t *testing.T) {
	const hintMessage = `
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "github.com/RedSock" and final result is "github.com/RedSock/rscli"
>`
	tmpDir := path.Join(os.TempDir(), "rscliTest"+strconv.Itoa(rand.Int()))

	initProjectF := func(t *testing.T) (tmpTestDir string, clean func()) {
		tmpTestDir = tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpTestDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		clean = func() {
			err = os.RemoveAll(tmpTestDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}
		defer func() {
			if err != nil {
				clean()
			}
		}()
		ioMock := mocks.NewIOMock(t)

		expectedPrint := []string{hintMessage}

		ioMock.PrintMock.Set(func(in string) {
			if in[0] == '\033' {
				return
			}

			if len(expectedPrint) == 0 {
				require.Failf(t, "unexpected message came in", "got %s", in)
			}
			require.Equal(t, expectedPrint[0], in)
			expectedPrint = expectedPrint[1:]
		})
		ioMock.GetInputMock.Expect().Return(projectName, nil)

		expectedColors := []struct {
			color colors.Color
			text  []string
		}{
			{
				color: colors.ColorCyan,
				text:  []string{`Wonderful!!! "github.com/RedSock/rscli" it is!`},
			},
			{
				color: colors.ColorGreen,
				text: []string{fmt.Sprintf(`Done.
Initialized new project github.com/RedSock/rscli
at %s/%s`, tmpTestDir, "rscli")},
			},
		}

		ioMock.PrintlnColoredMock.Set(func(color colors.Color, in ...string) {
			if len(expectedColors) == 0 || len(expectedColors[0].text) != len(in) {
				require.Failf(t, "unexpected message came in", "got %s with color %v", in, color)
			}

			require.Equal(t, expectedColors[0].color, color)
			for i, word := range in {
				require.Equal(t, expectedColors[i].text[i], word)
			}
			expectedColors = expectedColors[1:]
		})

		expectedPrintln := []string{
			"Starting project constructor",
			"_ ", "Preparing project structure",
			"_ ", "Preparing environment folder",
			"_ ", "Preparing config folder",
			"_ ", "Building project",
			"_ ", "Initiating go project",
			"_ ", "Cleaning up the project",
			"_ ", "Performing project fix up",
			"_ ", "Initiating git",
		}

		ioMock.PrintlnMock.Set(func(in ...string) {
			if len(expectedPrintln) == 0 || len(expectedPrintln) < len(in) {
				require.Failf(t, "unexpected message came in", "got %s", in)
			}
			for _, item := range in {
				require.Equal(t, expectedPrintln[0], item)
				expectedPrintln = expectedPrintln[1:]
			}
		})

		p := projectInit{
			io: ioMock,
			config: &config.RsCliConfig{
				Env: config.Project{
					PathToConfig: stdConfigPath,
				},
				DefaultProjectGitPath: "github.com/RedSock",
			},
			path: tmpTestDir,
		}

		cmd := newInitCmd(p)
		cmd.SetArgs([]string{""})
		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())

		return tmpTestDir, clean
	}

	testCases := []struct {
		name       string
		validators []projectAddValidator
	}{
		{
			name: "OK_ADD_ALL_DEPENDENCIES",
			validators: []projectAddValidator{
				redisValidator{},
				postgresValidator{},
				telegramValidator{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpTestDir, clean := initProjectF(t)
			defer clean()

			ioMock := mocks.NewIOMock(t)

			projAdd := projectAdd{
				io: ioMock,
			}

			projAdd.config = &config.RsCliConfig{
				Env: config.Project{
					PathToServers:  []string{stdServerPath},
					PathToConfig:   stdConfigPath,
					PathsToClients: []string{stdConfigPathToClient},
				},
			}

			cmd := newAddCmd(projAdd)
			err := cmd.Flags().Set(pathFlag, path.Join(tmpTestDir, projectName))
			require.NoError(t, err, "error setting path flag value")
			cmd.SetArgs([]string{resources.PostgresResourceName, resources.RedisResourceName, resources.TelegramResourceName})

			err = cmd.Execute()
			require.NoError(t, err, "error while initiating project")
			require.True(t, ioMock.MinimockPrintDone())
			require.True(t, ioMock.MinimockGetInputDone())
			require.True(t, ioMock.MinimockPrintlnColoredDone())
		})
	}

}

type redisValidator struct{}

func (r redisValidator) getName() string {
	return resources.RedisResourceName
}
func (r redisValidator) validate(t *testing.T, pth string) {
	pathToRedisConn := path.Join(pth, stdConfigPathToClient, resources.RedisResourceName, projpatterns.ConnFileName)
	f, err := os.ReadFile(pathToRedisConn)
	require.NoError(t, err, "error reading redis conn file")

	require.Equal(t, projpatterns.RedisConnFile.Content, f)
}

type postgresValidator struct{}

func (r postgresValidator) getName() string {
	return resources.PostgresResourceName
}
func (r postgresValidator) validate(t *testing.T, pth string) {
	pathToPGFolder := path.Join(pth, stdConfigPathToClient, resources.PostgresResourceName)

	f, err := os.ReadFile(path.Join(pathToPGFolder, projpatterns.PgConnFile.Name))
	require.NoError(t, err, "error reading postgres conn file")
	require.Equal(t, projpatterns.PgConnFile.Content, f)

	f, err = os.ReadFile(path.Join(pathToPGFolder, projpatterns.PgTxFileName))
	require.NoError(t, err, "error reading postgres tx file")
	require.Equal(t, projpatterns.PgTxFile.Content, f)
}

type telegramValidator struct{}

func (r telegramValidator) getName() string {
	return resources.TelegramResourceName
}
func (r telegramValidator) validate(t *testing.T, pth string) {
	pathToTgFolder := path.Join(pth, stdConfigPathToClient, resources.TelegramResourceName)

	f, err := os.ReadFile(path.Join(pathToTgFolder, projpatterns.ConnFileName))
	require.NoError(t, err, "error reading telegram conn file")
	require.Equal(t, projpatterns.TgConnFile, f)
}
