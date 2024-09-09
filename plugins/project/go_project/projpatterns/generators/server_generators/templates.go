package server_generators

import (
	_ "embed"
	"text/template"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/server"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

var (
	//go:embed templates/fs.go.pattern
	fileServerPattern  string
	fileServerTemplate *template.Template
)

func init() {
	fileServerTemplate = template.Must(
		template.New("config_autoload").
			Parse(fileServerPattern))
}

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
		return nil, errors.Wrap(err, "error executing file server template")
	}

	f := &folder.Folder{
		Name:    projpatterns.AppInitServerFileName,
		Content: file.Bytes(),
	}

	return f, nil
}
