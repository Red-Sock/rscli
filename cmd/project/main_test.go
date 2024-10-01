package project

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/tests/mocks"
)

const (
	testPath         = "test/"
	testExpectedPath = "test_expectation/"
)

func TestMain(m *testing.M) {
	err := rscliconfig.InitConfig(nil, nil)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func initNewProject(t *testing.T, name string) string {
	projInit := projectInit{
		io:     getEmptyIoMock(t),
		config: rscliconfig.GetConfig(),
		path:   testFolder,
	}

	folderPath := filepath.Join(testFolder, name)

	require.NoError(t, os.RemoveAll(folderPath))
	require.NoError(t, os.MkdirAll(folderPath, 0755))

	require.NoError(t, projInit.run(nil, []string{name}))

	return folderPath
}

func getEmptyIoMock(t *testing.T) *mocks.IOMock {
	ioMock := mocks.NewIOMock(t)
	ioMock.PrintlnColoredMock.Set(func(color colors.Color, in ...string) {})
	ioMock.PrintMock.Set(func(in string) {})
	ioMock.PrintlnMock.Set(func(in ...string) {})

	return ioMock
}

func compareDirs(t *testing.T, actualPath, expectedPath string) {
	dirs, err := os.ReadDir(expectedPath)
	require.NoError(t, err)

	for _, d := range dirs {
		if d.IsDir() {
			compareDirs(t, path.Join(actualPath, d.Name()), path.Join(expectedPath, d.Name()))
			continue
		}
		actualFilePath := path.Join(actualPath, d.Name())
		actualFile, err := os.ReadFile(actualFilePath)
		require.NoError(t, err)

		expectedFilePath := path.Join(expectedPath, d.Name())
		expectedFile, err := os.ReadFile(expectedFilePath)
		require.NoError(t, err)

		actualIdx, expectedIdx := 0, 0
		line := 1
		for actualIdx < len(actualFile) && expectedIdx < len(expectedFile) {
			if actualFile[actualIdx] == '\n' {
				line++
			}

			if actualFile[actualIdx] != expectedFile[expectedIdx] {
				wd, err := os.Getwd()
				require.NoError(t, err)

				t.Errorf("File         \n%s:%d\n is not equal to expected file \n%s:%d",
					wd+"/"+actualFilePath, line, wd+"/"+expectedFilePath, line)
				break
			}

			actualIdx++
			expectedIdx++
		}
	}
}

func getTestName(t *testing.T) string {
	name := t.Name()
	name = strings.ReplaceAll(name, "/", "_")[5:]
	name = strings.ToLower(name)
	return name
}
