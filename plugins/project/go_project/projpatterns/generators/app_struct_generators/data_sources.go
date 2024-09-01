package app_struct_generators

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies/grpc_discovery"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
)

func generateDataSourceInitFileAndArgs(dataSources matreshka.DataSources) (InitDepFuncGenArgs, []byte, error) {
	initDsArgs := InitDepFuncGenArgs{
		InitFunctionName: "InitDataSources",
		Imports:          make(map[string]string),
		Functions:        make([]InitFuncCall, 0, len(dataSources)),
	}

	for _, ds := range dataSources {
		fc := InitFuncCall{
			ResultName: generators.NormalizeResourceName(ds.GetName()),
		}

		var importPath, importAlias string
		switch ds.GetType() {
		case resources.PostgresResourceName, resources.SqliteResourceName:
			importPath, importAlias = sqlInitFunc(&fc)
		case resources.RedisResourceName:
			importPath, importAlias = redisInitFunc(&fc)
		case resources.TelegramResourceName:
			importPath, importAlias = telegramInitFunc(&fc)
		case resources.GrpcResourceName:
			var err error
			importPath, importAlias, err = grpcInitFunc(ds, &fc)
			if err != nil {
				return initDsArgs, nil, errors.Wrap(err, "error creating init func for grpc client")
			}
		default:
			return initDsArgs, nil, errors.New("unknown resource " + ds.GetType())
		}
		initDsArgs.Imports[importPath] = importAlias

		initDsArgs.Functions = append(initDsArgs.Functions, fc)
	}

	initDsArgs.Imports["github.com/Red-Sock/trace-errors"] = "errors"

	file := &rw.RW{}
	err := initAppStructFuncTemplate.Execute(file, initDsArgs)
	if err != nil {
		return initDsArgs, nil, errors.Wrap(err, "error generating server init file")
	}

	return initDsArgs, file.Bytes(), nil
}

func sqlInitFunc(fc *InitFuncCall) (importPath, importAlias string) {
	fc.FuncName = "sqldb.New"
	fc.ResultType = "*sqldb.DB"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during sql connection initialization"

	return "proj_name/internal/clients/sqldb", ""
}

func redisInitFunc(fc *InitFuncCall) (importPath, importAlias string) {
	fc.FuncName = "redis.New"
	fc.ResultType = "*redis.Client"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during redis connection initialization"

	return "proj_name/internal/clients/redis", ""
}

func telegramInitFunc(fc *InitFuncCall) (importPath, importAlias string) {
	fc.FuncName = "telegram.New"
	fc.ResultType = "*telegram.Bot"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during telegram bot initialization"

	return "proj_name/internal/clients/telegram", ""
}

func grpcInitFunc(grpc resources.Resource, fc *InitFuncCall) (importPath, importAlias string, err error) {
	grpcRes, ok := grpc.(*resources.GRPC)
	if !ok {
		return "", "", errors.New("not a grpc struct")
	}

	grpcPackage, err := grpc_discovery.DiscoverPackage(grpcRes.Module)
	if err != nil {
		return "", "", errors.Wrap(err, "error discovering grpc package")
	}

	fc.FuncName = "grpc." + grpcPackage.Constructor
	fc.ResultType = "grpc." + grpcPackage.ClientName
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during grpc client initialization"

	return "proj_name/internal/clients/grpc", "", nil
}
