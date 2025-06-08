package grpc_api

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project"
)

type serviceProtoApiArgs struct {
	PackageName    string
	GoPackageName  string
	NpmPackageName string
}

func GenerateServiceApiProto(project project.IProject) (*folder.Folder, error) {
	args := serviceProtoApiArgs{
		PackageName: project.GetShortName(),
	}

	protoFile := &rw.RW{}
	err := basicApiProtoTemplate.Execute(protoFile, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating service api proto")
	}

	return &folder.Folder{
		Name:    project.GetShortName() + ".proto",
		Content: protoFile.Bytes(),
	}, nil
}
