package app_struct_generators

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
)

func generateDataSourceInitFileAndArgs(dataSources matreshka.DataSources) (InitDepFuncGenArgs, []byte, error) {
	if len(dataSources) == 0 {
		return InitDepFuncGenArgs{}, nil, nil
	}

	initDsArgs := InitDepFuncGenArgs{
		InitFunctionName: "InitDataSources",
		Imports:          make(map[string]string),
		Functions:        make([]InitFuncCall, 0, len(dataSources)),
	}

	for _, ds := range dataSources {
		fc := InitFuncCall{
			ResultName: generators.NormalizeResourceName(ds.GetName()),
		}
		switch ds.GetType() {
		case resources.PostgresResourceName, resources.SqliteResourceName:
			fc.FuncName = "sqldb.New"
			fc.ResultType = "*sqldb.DB"
			initDsArgs.Imports["proj_name/internal/clients/sqldb"] = ""
			fc.ErrorMessage = "error during sql connection initialization"
		case resources.RedisResourceName:
			fc.FuncName = "redis.New"
			fc.ResultType = "*redis.Client"
			initDsArgs.Imports["proj_name/internal/clients/redis"] = ""
			fc.ErrorMessage = "error during redis connection initialization"
		case resources.TelegramResourceName:
			fc.FuncName = "telegram.New"
			fc.ResultType = "*telegram.Bot"
			initDsArgs.Imports["proj_name/internal/clients/telegram"] = ""
			fc.ErrorMessage = "error during telegram bot initialization"
		default:
			continue
		}

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
