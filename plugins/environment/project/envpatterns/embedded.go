package envpatterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

const (
	PortSuffix = "PORT"

	EnvDir = "environment"
)

const (
	ProjNamePattern            = "proj_name"
	ProjNameCapsPattern        = "PROJ_NAME_CAPS"
	AbsoluteProjectPathPattern = "abs_proj_path"
	PathToMain                 = "path_to_main"
	ResourceCapsPattern        = "RESOURCE"
	ResourceNameCapsPattern    = ResourceCapsPattern + "_NAME_CAPS"
)

const (
	HostEnvSuffix = "_HOST"
	Localhost     = "0.0.0.0"
)

const (
	MakefileEnvUpRuleName   = "env-up"
	MakefileEnvDownRuleName = "env-down"
)

var (
	//go:embed files/.env
	envFile []byte
	EnvFile = folder.Folder{
		Name:    ".env",
		Content: envFile,
	}
)

var (
	//go:embed files/docker-compose.yaml
	mainServiceComposeFile []byte
	DockerComposeFile      = folder.Folder{
		Name:    "docker-compose.yaml",
		Content: mainServiceComposeFile,
	}
)

var (
	//go:embed files/Makefile
	makefile []byte
	Makefile = folder.Folder{
		Name:    "Makefile",
		Content: makefile,
	}
)

var (
	//go:embed files/compose.examples.yaml
	buildInComposeExamples []byte
	BuildInComposeExamples = folder.Folder{
		Name:    "compose.examples.yaml",
		Content: buildInComposeExamples,
	}
)
