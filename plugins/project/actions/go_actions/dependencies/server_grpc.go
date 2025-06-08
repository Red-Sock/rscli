package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/grpc_api"
)

const grpcServerBasePath = "/{GRPC}"

type GrpcServer struct {
	dependencyBase
}

func grpcServer(dep dependencyBase) Dependency {
	return &GrpcServer{
		dep,
	}
}

func (r GrpcServer) GetFolderName() string {
	if r.Name != "" {
		return r.Name
	}

	return "grpc"
}

func (r GrpcServer) AppendToProject(proj Project) error {
	protoName := proj.GetShortName() + "_api.proto"

	ok, err := containsDependencyFolder(
		[]string{r.Cfg.Env.PathToServerDefinition},
		proj.GetFolder(),
		protoName)
	if err != nil {
		return rerrors.Wrap(err, "error searching dependencies")
	}

	if !ok {
		protoPath := path.Join(r.Cfg.Env.PathToServerDefinition, r.GetFolderName(), protoName)
		err = r.applyApiFolder(proj, protoPath)
		if err != nil {
			return rerrors.Wrap(err, "error applying grpc api folder")
		}
	}

	r.addGrpcServerToConfig(proj)

	initServerManagerFiles(proj)

	return nil
}

func (r GrpcServer) applyApiFolder(proj Project, protoPath string) error {
	protoFile, err := grpc_api.GenerateServiceApiProto(proj)
	if err != nil {
		return rerrors.Wrap(err)
	}

	protoFile.Name = protoPath

	proj.GetFolder().Add(protoFile)
	return nil
}

func (r GrpcServer) addGrpcServerToConfig(proj Project) {
	srv := prepareServerConfig(proj)
	if len(srv.GRPC) != 0 {
		return
	}

	srv.GRPC[grpcServerBasePath] = &server.GRPC{
		Module:  proj.GetName(),
		Gateway: "/api",
	}
}
