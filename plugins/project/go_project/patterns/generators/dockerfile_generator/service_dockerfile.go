package dockerfile_generator

import (
	_ "embed"
	"sort"
	"strings"
	"text/template"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/config"
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
	args := serviceDockerfileArgs{}

	cfg := proj.GetConfig()

	args.extractDataVolumes(cfg)
	args.extractPorts(cfg)

	out := rw.RW{}

	err := dockerfileTemplate.Execute(&out, args)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	return out.Bytes(), nil
}

func (args *serviceDockerfileArgs) extractDataVolumes(cfg *config.Config) {
	args.HasMigrations = false
	args.Volumes = nil

	for _, ds := range cfg.DataSources {
		switch v := ds.(type) {
		case *resources.Sqlite:
			args.HasMigrations = true
			args.Volumes = append(args.Volumes, v.Path)
		case *resources.Postgres:
			args.HasMigrations = true
		}
	}
}

func (args *serviceDockerfileArgs) extractPorts(cfg *config.Config) {
	ports := []string{}

	for _, srv := range cfg.Servers {
		ports = append(ports, srv.Port)
	}

	sort.Strings(ports)

	args.PortsList = strings.Join(ports, " ")
}
