package dependencies

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/api"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type Grpc struct {
	Cfg *rscliconfig.RsCliConfig
	Io  io.StdIO
}

func (r Grpc) GetFolderName() string {
	return "grpc"
}

func (r Grpc) AppendToProject(proj interfaces.Project) error {
	protoName := proj.GetShortName() + "_api.proto"

	ok, err := containsDependencyFolder(
		[]string{r.Cfg.Env.PathToServerDefinition},
		proj.GetFolder(),
		protoName)
	if err != nil {
		return errors.Wrap(err, "error searching dependencies")
	}

	protoPath := path.Join(r.Cfg.Env.PathToServerDefinition, r.GetFolderName(), protoName)

	if !ok {
		err := r.applyApiFolder(proj, protoPath)
		if err != nil {
			return errors.Wrap(err, "error applying grpc api folder")
		}
	}

	r.applyMakefile(proj)
	r.applyConfig(proj)

	return nil
}

func (r Grpc) applyApiFolder(proj interfaces.Project, protoPath string) error {
	serverF := projpatterns.ProtoServer.CopyWithNewName(protoPath)

	serverF.Content = renamer.ReplaceProjectNameShort(serverF.Content, proj.GetShortName())

	proj.GetFolder().Add(serverF)

	return nil
}

func (r Grpc) applyMakefile(proj interfaces.Project) {
	f := proj.GetFolder().GetByPath(projpatterns.GrpcMK.Name)
	if f != nil {
		return
	}

	proj.GetFolder().Add(projpatterns.GrpcMK.Copy())
}

func (r Grpc) applyConfig(proj interfaces.Project) {
	for _, item := range proj.GetConfig().Servers {
		if item.GetName() == r.GetFolderName() {
			return
		}
	}

	proj.GetConfig().Servers = append(proj.GetConfig().Servers,
		&api.Rest{
			Name: api.Name(r.GetFolderName()),
			Port: api.DefaultGrpcPort,
		})
}
