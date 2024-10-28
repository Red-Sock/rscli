package test_prepare

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/godverv/matreshka"
	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

var (
	//go:embed basic_config.yaml
	basicConfigFile []byte
)

type MockProject struct {
	Path string
	*project.Project

	rscliConfig *rscliconfig.RsCliConfig

	config matreshka.AppConfig
}

type opt func(m *MockProject)

func PrepareProject(t *testing.T, opts ...opt) *MockProject {
	p := &MockProject{
		rscliConfig: rscliconfig.GetConfig(),
	}

	p.config = matreshka.NewEmptyConfig()
	require.NoError(t, p.config.Unmarshal(basicConfigFile))

	if p.Path == "" {
		p.Path = t.Name()
	}

	for _, o := range opts {
		o(p)
	}

	require.NoError(t, os.MkdirAll(p.Path, 0777))

	cfgMarshalled, err := p.config.Marshal()
	require.NoError(t, err)

	p.WriteFile(t,
		path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile),
		cfgMarshalled,
	)

	p.Project, err = project.LoadProject(p.Path, p.rscliConfig)
	require.NoError(t, err)

	return p
}

func (m *MockProject) Clean(t *testing.T) {
	err := os.RemoveAll(m.Path)
	require.NoError(t, err)
}

func (m *MockProject) WriteFile(t *testing.T, relativePath string, data []byte) {
	relativePath = path.Join(m.Path, relativePath)

	require.NoError(t, os.MkdirAll(path.Dir(relativePath), 0777))

	cfgFile, err := os.Create(relativePath)
	require.NoError(t, err)

	_, err = cfgFile.Write(data)
	require.NoError(t, err)
	require.NoError(t, cfgFile.Close())
}
