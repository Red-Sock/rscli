package config_generators

import (
	"reflect"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators"
)

type structArgs struct {
	GenTag string

	StructName string
	Imports    map[string]string // path to alias
	Fields     []generators.KeyValue
}

func newConfigStructGenArgs(structName string) structArgs {
	return structArgs{
		GenTag:     defaultGenTag,
		StructName: structName,
		Imports:    map[string]string{},
		Fields:     nil,
	}
}

func getTypeName(in any) string {
	return reflect.ValueOf(in).Type().Name()
}
