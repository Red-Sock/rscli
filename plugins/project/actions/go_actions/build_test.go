package go_actions

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests"
	"github.com/Red-Sock/rscli/tests/project_mock"
)

type BuildProjectSuite struct {
	suite.Suite

	projectPath string
	expected    map[string][]byte
}

func (s *BuildProjectSuite) Test_BuildProject() {
	t := s.T()
	t.Parallel()

	mainGoFilePath := path.Join(patterns.CmdFolder, patterns.ServiceFolder, patterns.MainFile.Name)
	configPath := path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile)

	mainGoFile := patterns.MainFile.Copy().Content
	proj := project_mock.GetMockProject(t,
		project_mock.WithFile(mainGoFilePath, mainGoFile),
		project_mock.WithFileSystem(t),
	)

	s.projectPath = proj.Path

	action := BuildProjectAction{}

	require.NoError(t, action.Do(proj))

	dir, err := os.ReadDir(proj.Path)
	require.NoError(t, err)

	s.expected = map[string][]byte{
		mainGoFilePath: bytes.ReplaceAll(mainGoFile, []byte("proj_name"), []byte(proj.Name)),
		configPath:     project_mock.BasicConfig(),
	}

	s.walkDirAndValidate(s.projectPath, dir)
}

func (s *BuildProjectSuite) TearDownTest() {
	s.Require().NoError(os.RemoveAll(s.projectPath))
}

func (s *BuildProjectSuite) walkDirAndValidate(root string, dir []os.DirEntry) {
	for _, d := range dir {
		name := d.Name()

		if d.IsDir() {
			newRoot := path.Join(root, name)
			innerDirs, err := os.ReadDir(newRoot)
			s.Require().NoError(err)
			s.walkDirAndValidate(newRoot, innerDirs)
			continue
		}

		filePath := path.Join(root, name)
		expectedContent, ok := s.expected[filePath[len(s.projectPath)+1:]]
		if !ok {
			s.Assert().Fail("unexpected file name", name)
			continue
		}
		actualContent, err := os.ReadFile(filePath)
		s.Require().NoError(err)

		tests.CompareLongStrings(s.T(), expectedContent, actualContent)
	}
}

func Test_BuildProject(t *testing.T) {
	suite.Run(t, new(BuildProjectSuite))
}
