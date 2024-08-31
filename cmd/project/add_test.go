package project

import (
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
)

func Test_Add(t *testing.T) {
	t.Run("OK_ALL", func(t *testing.T) {

		projAdd := projectAdd{
			io:     getEmptyIoMock(t),
			config: rscliconfig.GetConfig(),
		}

		testName := getTestName(t)

		projAdd.path = initNewProject(t, testName)

		require.NoError(t, projAdd.run(nil, []string{
			"postgres",
			"redis",
			"grpc",
			"rest",
			"telegram",
			"sqlite",
		}))

		compareDirs(t, testPath+testName, testExpectedPath+testName)
	})
}
