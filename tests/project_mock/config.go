package project_mock

import (
	_ "embed"
	"time"

	"github.com/godverv/matreshka/environment"
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
		{
			Name:  "test_string_variable",
			Type:  environment.VariableTypeStr,
			Value: "test_value",
		},
		{
			Name:  "test_int_variable",
			Type:  environment.VariableTypeInt,
			Value: 1,
		},
		{
			Name:  "test_bool_variable",
			Type:  environment.VariableTypeBool,
			Value: true,
		},
		{
			Name:  "test_float_variable",
			Type:  environment.VariableTypeFloat,
			Value: 1.1,
		},
		{
			Name:  "test_duration_variable",
			Type:  environment.VariableTypeDuration,
			Value: time.Second * 5,
		},
		// Multiple values
		{
			Name:  "test_string_variables",
			Type:  environment.VariableTypeStr,
			Value: []string{"test_value", "test_value2"},
		},
		{
			Name:  "test_int_variables",
			Type:  environment.VariableTypeInt,
			Value: []int{1, 2},
		},
		{
			Name:  "test_bool_variables",
			Type:  environment.VariableTypeBool,
			Value: []bool{true, false},
		},
		{
			Name:  "test_float_variables",
			Type:  environment.VariableTypeFloat,
			Value: []float64{1.1, 2.2},
		},
		{
			Name:  "test_duration_variables",
			Type:  environment.VariableTypeDuration,
			Value: []time.Duration{time.Second * 5, time.Second * 8},
		},
	}
}
