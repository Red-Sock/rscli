package config_generators

import (
	"testing"
	"time"

	"github.com/godverv/matreshka"
	"github.com/stretchr/testify/require"
)

func Test_GenerateConfig(t *testing.T) {
	t.Run("environment", func(t *testing.T) {
		env := matreshka.Environment{
			{
				Name:  "one",
				Type:  "int",
				Value: 1,
			},
			{
				Name:  "two",
				Type:  "duration",
				Value: time.Second,
			},
		}

		_, generatedFolder, err := newGenerateEnvironmentConfigStruct(env)()
		require.NoError(t, err)

		expected := `// Code generated by RedSock CLI. DO NOT EDIT.

package config

import ( 
     "time"
)

type EnvironmentConfig struct { 
    One int
    Two time.Duration
}
`
		require.Equal(t, expected, string(generatedFolder.Content))
	})
}
