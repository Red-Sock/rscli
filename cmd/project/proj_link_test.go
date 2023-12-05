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

	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/tests/mocks"
)

func Test_LinkProject(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "rscliTest"+strconv.Itoa(rand.Int()))

	t.Run("OK_LINK_ONE_PROJECT", func(t *testing.T) {
		t.Parallel()

		tmpTestDir := tmpDir + "_" + strings.Split(t.Name(), "/")[1]
		defer func() {
			err := os.RemoveAll(tmpTestDir)
			require.NoError(t, err, "error during tmp dir deletion")
		}()

		ioMock := mocks.NewIOMock(t)

		expectedPrint := []string{hintInitMessage}

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
at %s`, tmpTestDir)},
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

		createProject(t, tmpTestDir, ioMock)

		ioMock = mocks.NewIOMock(t)
		pl := projectLink{}

		cmd := newLinkCmd(pl)

		err := cmd.Flags().Set(pathFlag, tmpTestDir)
		require.NoError(t, err, "error setting path flag value")

		cmd.SetArgs([]string{""})

		err = cmd.Execute()
		require.NoError(t, err, "error while initiating project")
		require.True(t, ioMock.MinimockPrintDone())
		require.True(t, ioMock.MinimockGetInputDone())
		require.True(t, ioMock.MinimockPrintlnColoredDone())
	})

}
