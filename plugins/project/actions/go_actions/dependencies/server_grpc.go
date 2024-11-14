package dependencies

import (
	"bytes"
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/gobeam/stringy"
	"github.com/godverv/matreshka/server"

	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

const grpcServerBasePath = "/{GRPC}"

type GrpcServer struct {
	dependencyBase
}

func (r GrpcServer) GetFolderName() string {
	if r.Name != "" {
		return r.Name
	}

	return "grpc_impl"
}

func (r GrpcServer) AppendToProject(proj Project) error {
	protoName := proj.GetShortName() + "_api.proto"

	ok, err := containsDependencyFolder(
		[]string{r.Cfg.Env.PathToServerDefinition},
		proj.GetFolder(),
		protoName)
	if err != nil {
		return errors.Wrap(err, "error searching dependencies")
	}

	if !ok {
		protoPath := path.Join(r.Cfg.Env.PathToServerDefinition, r.GetFolderName(), protoName)
		err := r.applyApiFolder(proj, protoPath)
		if err != nil {
			return errors.Wrap(err, "error applying grpc api folder")
		}
	}

	r.addGrpcServerToConfig(proj)

	initServerManagerFiles(proj)

	return nil
}

func (r GrpcServer) applyApiFolder(proj Project, protoPath string) error {
	serverF := patterns.ProtoServer.CopyWithNewName(protoPath)

	projName := stringy.New(proj.GetShortName())

	serverF.Content = bytes.Replace(serverF.Content,
		[]byte(envpatterns.ProjNamePattern),
		[]byte(projName.SnakeCase().ToLower()),
		1)

	serverF.Content = bytes.Replace(serverF.Content,
		[]byte(envpatterns.ProjNamePattern),
		[]byte(projName.CamelCase().Get()),
		1)

	proj.GetFolder().Add(serverF)

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
