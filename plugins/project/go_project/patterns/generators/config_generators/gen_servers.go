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

func newGenerateServerConfigStruct(srv matreshka.Servers) internalConfigGenerator {
	return func() (InternalConfig, *folder.Folder, error) {
		ic := InternalConfig{
			FieldName:    "Servers",
			StructName:   "ServersConfig",
			From:         getTypeName(matreshka.Servers{}),
			ErrorMessage: "Error parsing servers to config",
		}

		ecg := newConfigStructGenArgs(ic.StructName)

		for _, s := range srv {
			var fieldKV generators.KeyValue
			fieldKV.Key = matreshka.ServerName(s.Name)

			refVal := reflect.ValueOf(s)
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
			return InternalConfig{}, nil, rerrors.Wrap(err, "error executing server config struct template")
		}

		f := &folder.Folder{
			Name:    patterns.ConfigServersFileName,
			Content: buf.Bytes(),
		}

		return ic, f, nil
	}
}
