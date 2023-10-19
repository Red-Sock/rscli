package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/plugins/shared/file"
)

const (
	RsCliMkFileName      = "rscli.mk"
	DockerfileFileName   = "Dockerfile"
	ReadMeFileName       = "README.md"
	GitignoreFileName    = ".gitignore"
	GolangCIYamlFileName = ".golangci.yaml"
	EnvConfigYamlFile    = "env.yaml"
	DevConfigYamlFile    = "dev.yaml"
	ConfigYamlFile       = "config.yaml"
)

// Build and deploy
var (
	//go:embed pattern_c/Dockerfile
	dockerfile []byte
	Dockerfile = file.File{
		Name:    DockerfileFileName,
		Content: dockerfile,
	}

	//go:embed pattern_c/.gitignore
	GitIgnore []byte

	//go:embed pattern_c/.golangci.yaml
	Linter []byte
)

// Documentation
var (
	//go:embed pattern_c/README.md
	Readme []byte
)

// Testing files
var (
	//go:embed pattern_c/examples/api.http
	ApiHTTP []byte
)

var (
	//go:embed pattern_c/rscli.mk

	RscliMK []byte
)
