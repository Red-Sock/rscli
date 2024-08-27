package project

import (
	"testing"
)

const (
	testFolder = "test"
)

func Test_InitProject(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		_ = initNewProject(t, "init_ok")
		// TODO check every thing is generated ok
	})
}
