package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/tests/mocks"
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
