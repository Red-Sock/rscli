package app_struct_generators

import (
	"strconv"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
)

var (
	ErrServerMustHaveName = errors.New("application contains more than two servers. Names required to be specified")
)

func generateServerInitFileAndArgs(servers matreshka.Servers) (InitDepFuncGenArgs, []byte, error) {
	if len(servers) == 0 {
		return InitDepFuncGenArgs{}, nil, nil
	}

	genArgs := InitDepFuncGenArgs{
		Imports:          make(map[string]string),
		InitFunctionName: "InitServers",
	}

	serversMustHaveNames := len(servers) > 1

	for port, server := range servers {
		portStr := strconv.Itoa(port)
		if serversMustHaveNames && server.Name == "" {
			return InitDepFuncGenArgs{}, nil,
				errors.Wrap(ErrServerMustHaveName,
					"server with port "+portStr+" doesn't have name")
		}

		initFuncCall := InitFuncCall{
			FuncName:     "transport.NewServerManager",
			Args:         "a.Ctx, \":" + portStr + "\"",
			ResultName:   "Server" + generators.NormalizeResourceName(server.Name),
			ResultType:   "*transport.ServersManager",
			ErrorMessage: "error during server initialization on port: " + portStr,
		}
		if server.Name != "" {
			initFuncCall.ErrorMessage += ", with name: " + server.Name
		}

		genArgs.Functions = append(genArgs.Functions, initFuncCall)
	}

	genArgs.Imports["github.com/Red-Sock/trace-errors"] = "errors"
	genArgs.Imports["proj_name/internal/transport"] = ""

	initServer := &rw.RW{}
	err := initAppStructFuncTemplate.Execute(initServer, genArgs)
	if err != nil {
		return genArgs, nil, errors.Wrap(err, "error generating server init file")
	}

	return genArgs, initServer.Bytes(), nil
}
