package config_generators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.verv.tech/matreshka"
	"go.verv.tech/matreshka/environment"
)

func Test_GenerateEnvConfig(t *testing.T) {
	t.Run("environment", func(t *testing.T) {
		env := matreshka.Environment{
			environment.MustNewVariable("one", 1),
			environment.MustNewVariable("two", time.Second),
		}

		_, generatedFolder, err := newGenerateEnvironmentConfigStruct(env)()
		require.NoError(t, err)

		expected := `
// Code generated by RedSock CLI. DO NOT EDIT.

package config

import ( 
     "time"
)

type EnvironmentConfig struct { 
    One int
    Two time.Duration
}
`[1:]
		require.Equal(t, expected, string(generatedFolder.Content))
	})
}
