package dockerfile

import (
	_ "embed"
	"sort"
	"strings"
	"text/template"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project"
)

var (
	//go:embed templates/Dockerfile.pattern
	dockerfile         string
	dockerfileTemplate *template.Template
)

func init() {
	dockerfileTemplate = template.
		Must(template.New("Dockerfile").
			Parse(dockerfile))
}

type serviceDockerfileArgs struct {
	Volumes       []string
	PortsList     string
	HasMigrations bool
}

func GenerateDockerfile(proj project.IProject) ([]byte, error) {
	args := serviceDockerfileArgs{
		Volumes:       nil,
		PortsList:     "",
		HasMigrations: false,
	}

	cfg := proj.GetConfig()

	for _, ds := range cfg.DataSources {
		sqlite, ok := ds.(*resources.Sqlite)
		if !ok {
			continue
		}

		args.HasMigrations = true
		args.Volumes = append(args.Volumes, sqlite.Path)
	}

	ports := []string{}

	for _, srv := range cfg.Servers {
		ports = append(ports, srv.Port)
	}

	sort.Strings(ports)

	args.PortsList = strings.Join(ports, ", ")

	out := rw.RW{}

	err := dockerfileTemplate.Execute(&out, args)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	return out.Bytes(), nil
}
