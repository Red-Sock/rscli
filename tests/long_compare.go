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

func CompareLongStrings(t *testing.T, expected, actual []byte) (eq bool) {
	expectedReader := bytes.NewReader(expected)
	actualReader := bytes.NewReader(actual)

	for {
		expectedSlice := make([]byte, 800)
		actualSlice := make([]byte, 800)

		expLen, expErr := expectedReader.Read(expectedSlice)
		actLen, actErr := actualReader.Read(actualSlice)

		if expErr == io.EOF && actErr == io.EOF {
			break
		}

		expectedSlice = expectedSlice[:expLen]
		actualSlice = actualSlice[:actLen]

		eq = assert.Equal(t, string(expectedSlice), string(actualSlice))
		if !eq {
			return false
		}
		require.NoError(t, actErr)
		require.NoError(t, expErr)
	}

	return true
}

func AssertFolderInFs(t *testing.T, dirPath string, expected *folder.Folder) {
	if len(expected.Content) != 0 {
		targetPath := path.Join(dirPath, expected.Name)
		file, err := os.ReadFile(targetPath)
		require.NoError(t, err)
		eq := false
		if len(expected.Content) < 800 {
			eq = assert.Equal(t, string(expected.Content), string(file))
		} else {
			eq = CompareLongStrings(t, expected.Content, file)
		}

		if !eq {
			assert.Failf(t, "contents not equal", "expected content of file %s to be same as %s", targetPath, expected.Name)
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
