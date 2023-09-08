package project

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project/processor/validators"
	"github.com/Red-Sock/rscli/tests/mocks"
)

func Test_InitProject(t *testing.T) {
	const hintMessage = `
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "github.com/RedSock" and final result is "github.com/RedSock/rscli"
>`
	tmpDir := path.Join(os.TempDir(), "rscliTest"+strconv.Itoa(rand.Int()))
	const pName = "gitlab.ru/redsock/rscli"
	t.Run("OK_NAME_AND_PATH_VIA_FLAG", func(t *testing.T) {
		t.Parallel()

		tmpDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

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
		ioMock.GetInputMock.Expect().Return("rscli", nil)

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
at %s`, tmpDir)},
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
		}

		cmd := newInitCmd(p)

		err = cmd.Flags().Set(nameFlag, "")
		require.NoError(t, err, "error setting name flag value")

		err = cmd.Flags().Set(pathFlag, tmpDir)
		require.NoError(t, err, "error setting path flag value")

		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

	t.Run("OK_NAME_WITH_SHORT_URL", func(t *testing.T) {
		t.Parallel()
		tmpDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

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
		ioMock.GetInputMock.Expect().Return(pName, nil)

		expectedColors := []struct {
			color colors.Color
			text  []string
		}{
			{
				color: colors.ColorCyan,
				text:  []string{`Wonderful!!! "` + pName + `" it is!`},
			},
			{
				color: colors.ColorGreen,
				text: []string{fmt.Sprintf(`Done.
Initialized new project `+pName+`
at %s`, tmpDir)},
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
		}

		cmd := newInitCmd(p)

		err = cmd.Flags().Set(nameFlag, "")
		require.NoError(t, err, "error setting name flag value")

		err = cmd.Flags().Set(pathFlag, tmpDir)
		require.NoError(t, err, "error setting path flag value")

		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

	t.Run("OK_NAME_WITH_FULL_URL", func(t *testing.T) {
		t.Parallel()

		tmpDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

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
		ioMock.GetInputMock.Expect().Return("https://"+pName, nil)

		expectedColors := []struct {
			color colors.Color
			text  []string
		}{
			{
				color: colors.ColorCyan,
				text:  []string{`Wonderful!!! "` + pName + `" it is!`},
			},
			{
				color: colors.ColorGreen,
				text: []string{fmt.Sprintf(`Done.
Initialized new project `+pName+`
at %s`, tmpDir)},
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
		}

		cmd := newInitCmd(p)

		err = cmd.Flags().Set(nameFlag, "")
		require.NoError(t, err, "error setting name flag value")

		err = cmd.Flags().Set(pathFlag, tmpDir)
		require.NoError(t, err, "error setting path flag value")

		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

	t.Run("OK_NAME_WITH_SHORT_URL", func(t *testing.T) {
		t.Parallel()

		tmpDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

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
		ioMock.GetInputMock.Expect().Return(pName, nil)

		expectedColors := []struct {
			color colors.Color
			text  []string
		}{
			{
				color: colors.ColorCyan,
				text:  []string{`Wonderful!!! "` + pName + `" it is!`},
			},
			{
				color: colors.ColorGreen,
				text: []string{fmt.Sprintf(`Done.
Initialized new project `+pName+`
at %s`, tmpDir)},
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
		}

		cmd := newInitCmd(p)

		err = cmd.Flags().Set(nameFlag, "")
		require.NoError(t, err, "error setting name flag value")

		err = cmd.Flags().Set(pathFlag, tmpDir)
		require.NoError(t, err, "error setting path flag value")

		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

	t.Run("ERROR_EMPTY_NAME", func(t *testing.T) {
		t.Parallel()

		tmpDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		err := os.MkdirAll(tmpDir, 0777)
		require.NoError(t, err, "error creating tmp dir")

		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

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
		ioMock.GetInputMock.Expect().Return("", nil)

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
		}
		cmd := newInitCmd(p)

		err = cmd.Flags().Set(pathFlag, tmpDir)
		require.NoError(t, err, "error while setting path flag")

		err = cmd.Execute()
		require.Contains(t, err.Error(), emptyNameErr.Error())

		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})
	t.Run("ERROR_INVALID_NAME", func(t *testing.T) {
		t.Parallel()

		cfg := &config.RsCliConfig{
			DefaultProjectGitPath: "github.com/RedSock",
		}

		ioMock := mocks.NewIOMock(t)

		ioMock.PrintMock.Expect(hintMessage)
		ioMock.GetInputMock.Expect().Return("rscli$1", nil)

		ioMock.PrintlnColoredMock.Expect(colors.ColorCyan, `Wonderful!!! "gitlab.com/RedSock/rscli" it is!`)

		p := projectInit{
			io: ioMock,
		}
		cmd := newInitCmd(p)
		err := cmd.Flag("name").Value.Set("")
		require.NoError(t, err, "error setting flag value")

		err = cmd.Execute()
		require.Contains(t, err.Error(), validators.ErrInvalidNameErr.Error())
		ioMock.MinimockPrintlnInspect()
		ioMock.MinimockGetInputInspect()
	})
}
