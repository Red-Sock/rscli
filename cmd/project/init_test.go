package project

import (
	"testing"
)

const (
	testFolder = "test"
)

func Test_Init_Project(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		testName := getTestName(t)

		_ = initNewProject(t, testName)

		compareDirs(t, testPath+testName, testExpectedPath+testName)
	})
}
