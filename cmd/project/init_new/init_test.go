package init_new

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests/mocks"
)

func Test_InitProject(t *testing.T) {
	dirPath := path.Join(os.TempDir(), t.Name())

	require.NoError(t, os.RemoveAll(dirPath))
	require.NoError(t, os.MkdirAll(dirPath, 0777))
	defer func() {
		//require.NoError(t, os.RemoveAll(dirPath))
	}()

	projFullName := defaultGitPath + "/" + projName

	io := mocks.NewIOMock(t)
	coloredPrintlnOutputs := []struct {
		color colors.Color
		text  string
	}{
		{
			color: colors.ColorCyan,
			text:  fmt.Sprintf(ackProjectNameMessagePattern, projFullName),
		},
		{
			color: colors.ColorGreen,
			text:  fmt.Sprintf(newProjectInitMessage, projFullName, dirPath+"/"+projName),
		},
	}
	coloredPrintlnIdx := 0

	io.PrintlnColoredMock.Set(func(color colors.Color, in ...string) {
		for _, text := range in {
			require.Equal(t, coloredPrintlnOutputs[coloredPrintlnIdx].color, color)
			require.Equal(t, coloredPrintlnOutputs[coloredPrintlnIdx].text, text)
			coloredPrintlnIdx++
		}
	})

	printlnOutputs := []string{
		"Starting project constructor",
		"Project actions performed",
	}
	printlnIndex := 0
	io.PrintlnMock.Set(func(in ...string) {
		for _, text := range in {
			require.Equal(t, printlnOutputs[printlnIndex], text)
			printlnIndex++
		}
	})

	args := []string{projName}

	cfg := &config.RsCliConfig{
		DefaultProjectGitPath: defaultGitPath,
	}

	basicProc := processor.New(
		processor.WithIo(io),
		processor.WithWd(dirPath),
		processor.WithConfig(cfg),
	)

	cmd := NewCommand(basicProc)

	cmd.SetArgs(args)
	err := cmd.Execute()
	require.NoError(t, err)

	projectPath := path.Join(dirPath, projName)

	assertFile(t, projectPath, patterns.Dockerfile)
	assertFile(t, projectPath, patterns.Makefile)

	readme := patterns.Readme.Copy()
	renamer.ReplaceProjectName(projFullName, readme)
	assertFile(t, projectPath, readme)

	rscliMk := patterns.RscliMK.Copy()
	renamer.ReplaceProjectName(projName, rscliMk)
	assertFile(t, projectPath, rscliMk)

}

func assertFile(t *testing.T, dirPath string, expectedFile *folder.Folder) {
	file, err := os.ReadFile(path.Join(dirPath, expectedFile.Name))
	require.NoError(t, err)
	assert.Equal(t, string(expectedFile.Content), string(file))
}
