package config_generators

import (
	"reflect"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
)

func newGenerateEnvironmentConfigStruct(environment matreshka.Environment) internalConfigGenerator {
	return func() (InternalConfig, *folder.Folder, error) {
		ic := InternalConfig{
			FieldName:    "Environment",
			StructName:   "EnvironmentConfig",
			From:         getTypeName(matreshka.Environment{}),
			ErrorMessage: "error parsing environment config",
		}

		ecg := newConfigStructGenArgs(ic.StructName)

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
			return InternalConfig{}, nil, errors.Wrap(err, "error executing config struct template")
		}

		f := &folder.Folder{
			Name:    projpatterns.ConfigEnvironmentFileName,
			Content: buf.Bytes(),
		}

		return ic, f, nil
	}
}
