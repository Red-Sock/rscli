package project

import (
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
)

func Test_Link(t *testing.T) {
	t.Run("GRPC", func(t *testing.T) {
		projLink := projectLink{
			io:     getEmptyIoMock(t),
			config: rscliconfig.GetConfig(),
		}
		testName := getTestName(t)
		projLink.path = initNewProject(t, testName)

		require.NoError(t,
			projLink.run(nil,
				[]string{"github.com/godverv/hello_world"}),
		)
		compareDirs(t, testPath+testName, testExpectedPath+testName)
	})
}
