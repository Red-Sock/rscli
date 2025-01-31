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

		generateReq := newConfigStructGenArgs(ic.StructName)

		for _, env := range environment {
			var fieldKV generators.KeyValue
			fieldKV.Key = generators.NormalizeResourceName(env.Name)

			v := env.Value.Value()

			if v != nil {
				refVal := reflect.ValueOf(v)
				tp := refVal.Type()
				fieldKV.Value = tp.String()

				if tp.PkgPath() != "" {
					generateReq.Imports[tp.PkgPath()] = ""
				}
			} else {
				fieldKV.Value = ""
			}

			enumVal := env.Enum.Value()
			if !env.Enum.IsZero() {
				switch vals := enumVal.(type) {
				case []string:
					enumToGen := EnumGenArg{
						Name: generators.NormalizeResourceName(env.Name),
					}

					enumToGen.Values = make([]generators.KeyValue, 0, len(vals))
					for _, old := range vals {
						enumToGen.Values = append(enumToGen.Values,
							generators.KeyValue{
								Key:   generators.NormalizeResourceName(old),
								Value: "\"" + old + "\"",
							})
					}
					generateReq.Enums = append(generateReq.Enums, enumToGen)
				case []int:
				default:
					return InternalConfig{}, nil,
						rerrors.New("error generating enums for config value. Unsupported enum type %T. Expected String slice",
							enumVal)
				}
			}

			generateReq.Fields = append(generateReq.Fields, fieldKV)
		}

		buf := &rw.RW{}
		err := configStructTemplate.Execute(buf, generateReq)
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
