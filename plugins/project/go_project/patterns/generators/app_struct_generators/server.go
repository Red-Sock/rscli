package app_struct_generators

import (
	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

var (
	ErrServerMustHaveName = rerrors.New("application contains more than two servers. Names required to be specified")
)

func generateServerInitFileAndArgs(servers matreshka.Servers) (InitServerListenersArgs, []byte, error) {
	genArgs := InitServerListenersArgs{
		Imports: make(map[string]string),
		Servers: nil,
	}

	serversMustHaveNames := len(servers) > 1

	for _, server := range servers {

		if serversMustHaveNames && server.Name == "" {
			return InitServerListenersArgs{}, nil,
				rerrors.Wrap(ErrServerMustHaveName,
					"server \""+server.Name+"\" doesn't exist in config")
		}

		name := matreshka.ServerName(server.Name)

		initFuncCall := InitServerListenerArgs{
			ServerName: name,
		}

		genArgs.Servers = append(genArgs.Servers, initFuncCall)
	}

	genArgs.Imports[patterns.ImportNameErrorsPackage] = ""
	genArgs.Imports["net"] = ""

	initServer := &rw.RW{}
	err := initServerTemplate.Execute(initServer, genArgs)
	if err != nil {
		return genArgs, nil, rerrors.Wrap(err, "error generating server init file")
	}

	return genArgs, initServer.Bytes(), nil
}
