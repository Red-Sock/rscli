package go_actions

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests/project_mock"
)

func Test_PrepareConfig(t *testing.T) {
	t.Parallel()

	type test struct {
		opts          []project_mock.Opt
		expectedFiles map[string][]byte
	}

	type testCase struct {
		new func() test
	}

	masterConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile)
	templateConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigTemplateYaml)
	devConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigDevYamlFile)

	tests := map[string]testCase{
		"generate_new_configs": {
			new: func() (tst test) {
				tst.opts = append(tst.opts,
					project_mock.WithEnvironmentVariables(
						project_mock.GetAllEnvVariables()...))

				fullConfig := project_mock.FullEnvConfig()

				tst.expectedFiles = map[string][]byte{
					masterConfigPath:   fullConfig,
					templateConfigPath: fullConfig,
					devConfigPath:      fullConfig,

					path.Join(patterns.InternalFolder, patterns.ConfigsFolder,
						patterns.ConfigEnvironmentFileName): fullConfigGoFile()}
				return tst
			},
		},
		"append_from_master_config": {
			new: func() (tst test) {
				// Test validates that if user added env var to master config
				// this value will be added to template and dev configs

				basicConfigFile := project_mock.BasicConfig()
				tst.opts = append(tst.opts,
					project_mock.WithFile(templateConfigPath, basicConfigFile),
					project_mock.WithFile(devConfigPath, basicConfigFile),
				)

				cfg := matreshka.NewEmptyConfig()
				require.NoError(t, cfg.Unmarshal(basicConfigFile))

				newEnvVariable := environment.MustNewVariable("test_value", "test")

				cfg.Environment = append(cfg.Environment, newEnvVariable)

				tst.opts = append(tst.opts,
					project_mock.WithEnvironmentVariables(newEnvVariable))

				masterConfigFile, err := cfg.Marshal()
				require.NoError(t, err)

				tst.opts = append(tst.opts,
					project_mock.WithFile(masterConfigPath, masterConfigFile),
				)

				tst.expectedFiles = map[string][]byte{
					masterConfigPath:   masterConfigFile,
					templateConfigPath: masterConfigFile,
					devConfigPath:      masterConfigFile,

					path.Join(patterns.InternalFolder, patterns.ConfigsFolder,
						patterns.ConfigEnvironmentFileName): testValueGoConfig(),
				}
				return tst
			},
		},
	}

	action := PrepareConfigFolder{}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc := tc.new()
			projectMock := project_mock.GetMockProject(t, tc.opts...)

			require.NoError(t, action.Do(projectMock.Project))

			for pathToFile, expected := range tc.expectedFiles {
				actualFile := projectMock.GetFolder().GetByPath(pathToFile)

				if strings.HasSuffix(pathToFile, ".yaml") {
					if !assert.YAMLEq(t, string(expected), string(actualFile.Content)) {
						require.YAMLEq(t, string(expected), string(actualFile.Content))
					}
				} else if string(expected) != string(actualFile.Content) {
					assert.Fail(t, pathToFile+" is not as expected")
					require.Equal(t, string(expected), string(actualFile.Content))
				}
			}
		})
	}
}

func fullConfigGoFile() []byte {
	return []byte(`
// Code generated by RedSock CLI. DO NOT EDIT.

package config

import ( 
     "time"
)

type EnvironmentConfig struct { 
    TestStringVariable string
    TestIntVariable int
    TestBoolVariable bool
    TestFloatVariable float64
    TestDurationVariable time.Duration
    TestStringVariables []string
    TestIntVariables []int
    TestBoolVariables []bool
    TestFloatVariables []float64
    TestDurationVariables []time.Duration
}
`)[1:]
}

func testValueGoConfig() []byte {
	return []byte(`
// Code generated by RedSock CLI. DO NOT EDIT.

package config

type EnvironmentConfig struct { 
    TestValue string
}
`)[1:]
}
