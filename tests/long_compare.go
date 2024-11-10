package tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
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
