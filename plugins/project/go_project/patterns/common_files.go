package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

const (
	EnvConfigYamlFile    = "env.yaml"
	ConfigDevYamlFile    = "dev.yaml"
	ConfigMasterYamlFile = "config.yaml"
	MakefileFile         = "Makefile"
	RscliMakefileFile    = "rscli.mk"
	DockerfileFile       = "Dockerfile"

	GenCommand           = "gen"
	GenGrpcServerCommand = "gen-server-grpc"
)

// Build and deploy
var (
	//go:embed pattern_c/.gitignore
	gitIgnore []byte
	GitIgnore = &folder.Folder{
		Name:    ".gitignore",
		Content: gitIgnore,
	}

	//go:embed pattern_c/.golangci.yaml
	linter []byte
	Linter = &folder.Folder{
		Name:    ".golangci.yaml",
		Content: linter,
	}
)

// Documentation
var (
	//go:embed pattern_c/README.md
	readme []byte
	Readme = &folder.Folder{
		Name:    "README.md",
		Content: readme,
	}
)

// Example files
var (
	//go:embed pattern_c/examples/api.http
	apiHTTP []byte
	ApiHTTP = &folder.Folder{
		Name:    "api.http",
		Content: apiHTTP,
	}
)

// Scripts
var (
	//go:embed pattern_c/rscli.mk
	rscliMK []byte
	RscliMK = &folder.Folder{
		Name:    "rscli.mk",
		Content: rscliMK,
	}

	//go:embed pattern_c/Makefile
	makefile []byte
	Makefile = &folder.Folder{
		Name:    MakefileFile,
		Content: makefile,
	}
)
