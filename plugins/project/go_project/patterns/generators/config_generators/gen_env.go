package config_generators

import (
	"reflect"

	"go.redsock.ru/rerrors"
	"go.verv.tech/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators"
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

			v := env.Value.Value()

			if v != nil {
				refVal := reflect.ValueOf(v)
				tp := refVal.Type()
				fieldKV.Value = tp.String()

				if tp.PkgPath() != "" {
					ecg.Imports[tp.PkgPath()] = ""
				}
			} else {
				fieldKV.Value = ""
			}

			ecg.Fields = append(ecg.Fields, fieldKV)
		}

		buf := &rw.RW{}
		err := configStructTemplate.Execute(buf, ecg)
		if err != nil {
			return InternalConfig{}, nil, rerrors.Wrap(err, "error executing config struct template")
		}

		f := &folder.Folder{
			Name:    patterns.ConfigEnvironmentFileName,
			Content: buf.Bytes(),
		}

		return ic, f, nil
	}
}
