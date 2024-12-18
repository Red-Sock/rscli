package app_struct_generators

import (
	"github.com/godverv/matreshka"
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators"
)

var (
	ErrServerMustHaveName = rerrors.New("application contains more than two servers. Names required to be specified")
)

func generateServerInitFileAndArgs(servers matreshka.Servers) (InitDepFuncGenArgs, []byte, error) {
	genArgs := InitDepFuncGenArgs{
		Imports:          make(map[string]string),
		InitFunctionName: "InitServers",
	}

	serversMustHaveNames := len(servers) > 1

	for _, server := range servers {

		if serversMustHaveNames && server.Name == "" {
			return InitDepFuncGenArgs{}, nil,
				rerrors.Wrap(ErrServerMustHaveName,
					"server \""+server.Name+"\" doesn't exist in config")
		}

		name := matreshka.ServerName(server.Name)

		initFuncCall := InitFuncCall{
			FuncName:     "transport.NewServerManager",
			Args:         "a.Ctx, a.Cfg.Servers." + name + ".Port",
			ResultName:   "Server" + generators.NormalizeResourceName(server.Name),
			ResultType:   "*transport.ServersManager",
			ErrorMessage: "error during \\\"" + matreshka.ServerName(server.Name) + "\\\" server initialization",
		}
		if server.Name != "" {
			initFuncCall.ErrorMessage += ", with name: " + server.Name
		}

		genArgs.Functions = append(genArgs.Functions, initFuncCall)
		genArgs.ServerName = initFuncCall.ResultName
	}

	genArgs.Imports[patterns.ImportNameErrorsPackage] = ""
	genArgs.Imports["proj_name/internal/transport"] = ""

	initServer := &rw.RW{}
	err := initAppStructFuncTemplate.Execute(initServer, genArgs)
	if err != nil {
		return genArgs, nil, rerrors.Wrap(err, "error generating server init file")
	}

	return genArgs, initServer.Bytes(), nil
}
