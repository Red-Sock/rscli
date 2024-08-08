package dependencies

import (
	"bytes"
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/servers"

	"github.com/gobeam/stringy"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type GrpcServer struct {
	Name string

	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (r GrpcServer) GetFolderName() string {
	if r.Name != "" {
		return r.Name
	}

	return "grpc"
}

func (r GrpcServer) AppendToProject(proj proj_interfaces.Project) error {
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

	r.applyConfig(proj)
	r.applyServerFolder(proj)

	applyServerFolder(proj)
	return nil
}

func (r GrpcServer) applyApiFolder(proj proj_interfaces.Project, protoPath string) error {
	serverF := projpatterns.ProtoServer.CopyWithNewName(protoPath)

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

func (r GrpcServer) applyConfig(proj proj_interfaces.Project) {
	for _, item := range proj.GetConfig().Servers {
		if item.GetName() == r.GetFolderName() {
			return
		}
	}

	proj.GetConfig().Servers = append(proj.GetConfig().Servers,
		&servers.Rest{
			Name: servers.Name(r.GetFolderName()),
			Port: servers.DefaultGrpcPort,
		})
}

func (r GrpcServer) applyServerFolder(proj proj_interfaces.Project) {
	f := proj.GetFolder()

	pth := []string{projpatterns.InternalFolder, projpatterns.TransportFolder, r.GetFolderName()}
	serverFolder := f.GetByPath(pth...)
	if serverFolder == nil {
		serverFolder = &folder.Folder{
			Name: path.Join(pth...),
		}
		f.Add(serverFolder)
	}

	if serverFolder.GetByPath(projpatterns.GrpcServFile.Name) == nil {
		serverFolder.Add(projpatterns.GrpcServFile.Copy())
	}
	// TODO генерация ручек-реализаций под конкракты
}
