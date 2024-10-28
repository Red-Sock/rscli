package project_mock

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/godverv/matreshka"
	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

type MockProject struct {
	*project.Project

	rscliConfig *rscliconfig.RsCliConfig
}

type Opt func(m *MockProject)

func GetMockProject(t *testing.T, opts ...Opt) *MockProject {
	p := &MockProject{
		rscliConfig: rscliconfig.GetConfig(),
		Project: &project.Project{
			Name: t.Name(),
			Cfg: &config.Config{
				AppConfig: matreshka.NewEmptyConfig(),
			},
			Root: folder.Folder{},
		},
	}

	require.NoError(t, p.Cfg.AppConfig.Unmarshal(basicConfigFile))

	for _, o := range opts {
		o(p)
	}

	cfgMarshalled, err := p.Cfg.AppConfig.Marshal()
	require.NoError(t, err)

	p.Cfg.AppConfig = matreshka.NewEmptyConfig()
	// to be shure in types of env variables
	require.NoError(t, p.Cfg.AppConfig.Unmarshal(cfgMarshalled))

	masterConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile)
	if p.Root.GetByPath(masterConfigPath) == nil {
		p.Root.Add(
			&folder.Folder{
				Name:    masterConfigPath,
				Content: cfgMarshalled,
			},
		)
	}

	return p
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
