package config_generators

import (
	"reflect"

	"go.redsock.ru/rerrors"
	"go.verv.tech/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies/link_service/grpc_discovery"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators"
)

func newGenerateDataSourcesConfigStruct(dataSources matreshka.DataSources) internalConfigGenerator {
	return func() (InternalConfig, *folder.Folder, error) {
		ic := InternalConfig{
			FieldName:    "DataSources",
			StructName:   "DataSourcesConfig",
			From:         getTypeName(matreshka.DataSources{}),
			ErrorMessage: "error parsing data sources to struct",
		}

		ecg := newConfigStructGenArgs(ic.StructName)

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
			return InternalConfig{}, nil, rerrors.Wrap(err, "error executing data source config struct template")
		}

		f := &folder.Folder{
			Name:    patterns.ConfigDataSourcesFileName,
			Content: buf.Bytes(),
		}

		return ic, f, nil
	}
}

func GenerateGRPCClient(args grpc_discovery.GrpcPackage) ([]byte, error) {
	buf := &rw.RW{}
	err := grpcConnectionTemplate.Execute(buf, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error executing template on grpc connection")
	}

	return buf.Bytes(), nil
}
