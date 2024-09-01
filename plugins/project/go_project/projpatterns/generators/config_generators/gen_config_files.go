package config_generators

import (
	_ "embed"
	"reflect"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
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

func GenerateEnvironmentConfigStruct(environment matreshka.Environment) ([]byte, error) {
	ecg := newConfigStructGenArgs("EnvironmentConfig")

	for _, env := range environment {
		var fieldKV generators.KeyValue
		fieldKV.Key = generators.NormalizeResourceName(env.Name)

		refVal := reflect.ValueOf(env.Value)
		tp := refVal.Type()

		fieldKV.Value = tp.String()

		if tp.PkgPath() != "" {
			ecg.Imports[tp.PkgPath()] = "" // todo think about aliases?
		}

		ecg.Fields = append(ecg.Fields, fieldKV)

	}

	buf := &rw.RW{}
	err := configStructTemplate.Execute(buf, ecg)
	if err != nil {
		return nil, errors.Wrap(err, "error executing template")
	}

	return buf.Bytes(), nil
}

func GenerateDataSourcesConfigStruct(dataSources matreshka.DataSources) ([]byte, error) {
	ecg := newConfigStructGenArgs("DataSourcesConfig")

	for _, ds := range dataSources {
		var fieldKV generators.KeyValue
		fieldKV.Key = generators.NormalizeResourceName(ds.GetName())

		refVal := reflect.ValueOf(ds)
		tp := refVal.Type()

		fieldKV.Value = tp.String()

		kind := tp.Kind()
		if kind == reflect.Ptr {
			tp = tp.Elem()
		}

		if tp.PkgPath() != "" {
			ecg.Imports[tp.PkgPath()] = "" // todo think about aliases?
		}

		ecg.Fields = append(ecg.Fields, fieldKV)

	}

	buf := &rw.RW{}
	err := configStructTemplate.Execute(buf, ecg)
	if err != nil {
		return nil, errors.Wrap(err, "error executing template")
	}

	return buf.Bytes(), nil
}

type GrpcClientArgs struct {
	ApiPackage  string
	Constructor string
	ClientName  string
}

func GenerateGRPCClient(args GrpcClientArgs) ([]byte, error) {
	buf := &rw.RW{}
	err := grpcConnectionTemplate.Execute(buf, args)
	if err != nil {
		return nil, errors.Wrap(err, "error executing template")
	}

	return buf.Bytes(), nil
}
