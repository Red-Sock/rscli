package init_new

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests"
	"github.com/Red-Sock/rscli/tests/mocks"
)

//go:embed expected/load_config/load_config_file.go
var loadConfigFile []byte

//go:embed expected/load_config/environment.go
var envConfigFile []byte

func Test_InitProject(t *testing.T) {
	dirPath := path.Join(os.TempDir(), t.Name())

	require.NoError(t, os.RemoveAll(dirPath))
	require.NoError(t, os.MkdirAll(dirPath, 0777))
	defer func() {
		require.NoError(t, os.RemoveAll(dirPath))
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

	tests.AssertFolderInFs(t, projectPath, patterns.Dockerfile)
	tests.AssertFolderInFs(t, projectPath, patterns.Makefile)

	readme := patterns.Readme.Copy()
	renamer.ReplaceProjectName(projFullName, readme)
	tests.AssertFolderInFs(t, projectPath, readme)

	rscliMk := patterns.RscliMK.Copy()
	renamer.ReplaceProjectName(projName, rscliMk)
	tests.AssertFolderInFs(t, projectPath, rscliMk)

	stat, err := os.Stat(path.Join(projectPath, patterns.GoMod))
	require.NoError(t, err)
	require.False(t, stat.IsDir())
	{
		basicConfig := []byte(`
app_info:
    name: github.com/RedSock/test_proj
    version: v0.0.1
    startup_duration: 10s
environment:
    - enum:
        - Trace
        - Debug
        - Info
        - Warn
        - Error
        - Fatal
        - Panic
      name: log-level
      type: string
      value: Info
    - enum:
        - JSON
        - TEXT
      name: log-format
      type: string
      value: TEXT
`)[1:]

		tests.AssertFolderInFs(t, projectPath,
			&folder.Folder{
				Name: patterns.ConfigsFolder,
				Inner: []*folder.Folder{
					{
						Name:    patterns.ConfigDevYamlFile,
						Content: basicConfig,
					},
					{
						Name:    patterns.ConfigMasterYamlFile,
						Content: basicConfig,
					},
					{
						Name:    patterns.ConfigTemplateYaml,
						Content: basicConfig,
					},
				},
			})
	}

	{
		tests.AssertFolderInFs(t, projectPath, &folder.Folder{
			Name: patterns.InternalFolder,
			Inner: []*folder.Folder{
				{
					Name: patterns.ConfigsFolder,
					Inner: []*folder.Folder{
						{
							Name:    patterns.ConfigLoadFileName,
							Content: loadConfigFile,
						},
						{
							Name:    patterns.ConfigEnvironmentFileName,
							Content: envConfigFile,
						},
					},
				},
			},
		})
	}

	mainGoFile := patterns.MainFile.Copy()
	renamer.ReplaceProjectName(projFullName, mainGoFile)
	tests.AssertFolderInFs(t, path.Join(projectPath, patterns.CmdFolder, patterns.ServiceFolder), mainGoFile)
}
