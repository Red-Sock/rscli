package server_generators

import (
	_ "embed"
	"text/template"

	"github.com/godverv/matreshka/server"
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

var (
	//go:embed templates/fs.go.pattern
	fileServerPattern  string
	fileServerTemplate = template.Must(template.New("file_server").Parse(fileServerPattern))
)

type fileServerGenArgs struct {
	DistPath string
}

func GenerateFileServer(fs server.FS) (*folder.Folder, error) {
	args := fileServerGenArgs{
		DistPath: fs.Dist,
	}

	file := &rw.RW{}
	err := fileServerTemplate.Execute(file, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error executing file server template")
	}

	f := &folder.Folder{
		Name:    patterns.AppInitServerFileName,
		Content: file.Bytes(),
	}

	return f, nil
}
