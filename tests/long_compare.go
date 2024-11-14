package tests

import (
	"bytes"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
)

func CompareLongStrings(t *testing.T, expected, actual []byte) {
	expectedReader := bytes.NewReader(expected)
	actualReader := bytes.NewReader(actual)
	for {
		expectedSlice := make([]byte, 800)
		actualSlice := make([]byte, 800)

		expLen, expErr := expectedReader.Read(expectedSlice)
		actLen, actErr := actualReader.Read(actualSlice)

		if expErr == io.EOF && actErr == io.EOF {
			return
		}

		expectedSlice = expectedSlice[:expLen]
		actualSlice = actualSlice[:actLen]

		require.Equal(t, string(expectedSlice), string(actualSlice))

		require.NoError(t, actErr)
		require.NoError(t, expErr)
	}
}

func AssertFolderInFs(t *testing.T, dirPath string, expected *folder.Folder) {
	if len(expected.Content) != 0 {
		file, err := os.ReadFile(path.Join(dirPath, expected.Name))
		require.NoError(t, err)
		if len(expected.Content) < 800 {
			assert.Equal(t, string(expected.Content), string(file))
		} else {
			CompareLongStrings(t, expected.Content, file)
		}
		return
	}

	for _, innerF := range expected.Inner {
		AssertFolderInFs(t, path.Join(dirPath, expected.Name), innerF)
	}
}

func AssertVirtualFolder(t *testing.T, proj project.IProject, expected *folder.Folder) {
	if len(expected.Content) != 0 {
		fileInProject := proj.GetFolder().GetByPath(expected.Name)
		require.NotNil(t, fileInProject, "file not found in project %s", expected.Name)

		if len(expected.Content) < 800 {
			assert.Equal(t, string(expected.Content), string(fileInProject.Content))
		} else {
			CompareLongStrings(t, expected.Content, fileInProject.Content)
		}
		return
	}

	for _, innerF := range expected.Inner {
		AssertVirtualFolder(t, proj, innerF)
	}
}
