package project_mock

import (
	_ "embed"
	"time"

	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
)

//go:embed full_config.yaml
var fullEnvConfigFile []byte

func FullEnvConfig() []byte {
	n := make([]byte, len(fullEnvConfigFile))
	copy(n, fullEnvConfigFile)

	return n
}

var (
	//go:embed basic_config.yaml
	basicConfigFile []byte
)

func BasicConfig() []byte {
	n := make([]byte, len(basicConfigFile))
	copy(n, basicConfigFile)

	return n
}

func GetAllEnvVariables() []*environment.Variable {
	return []*environment.Variable{
		// Single values
		environment.MustNewVariable("test_string_variable", "test_value"),
		environment.MustNewVariable("test_int_variable", 1),
		environment.MustNewVariable("test_bool_variable", true),
		environment.MustNewVariable("test_float_variable", 1.1),
		environment.MustNewVariable("test_duration_variable", time.Second*5),
		// Multiple values
		environment.MustNewVariable("test_string_variables", []string{"test_value", "test_value2"}),
		environment.MustNewVariable("test_int_variables", []int{1, 2}),
		environment.MustNewVariable("test_bool_variables", []bool{true, false}),
		environment.MustNewVariable("test_float_variables", []float64{1.1, 2.2}),
		environment.MustNewVariable("test_duration_variables", []time.Duration{time.Second * 5, time.Second * 8}),
	}
}
