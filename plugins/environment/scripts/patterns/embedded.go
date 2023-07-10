package patterns

import (
	_ "embed"
)

type File struct {
	Name    string
	Content []byte
}

const (
	EnvDir         = "environment"
	EnvExampleFile = ".env.example"
)

const (
	ProjNamePattern     = "proj_name"
	ProjNameCapsPattern = "PROJ_NAME_CAPS"
)

var (
	//go:embed files/.env
	envFile []byte
	EnvFile = File{
		Name:    ".env",
		Content: envFile,
	}
)

var (
	//go:embed files/docker-compose.yaml
	mainServiceComposeFile []byte
	DockerComposeFile      = File{
		Name:    "docker-compose.yaml",
		Content: mainServiceComposeFile,
	}
)

var (
	//go:embed files/Makefile
	makefile []byte
	Makefile = File{
		Name:    "Makefile",
		Content: makefile,
	}
)

//go:embed files/compose.examples.yaml
var buildInExamples []byte
