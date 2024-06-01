package projpatterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

const (
	EnvConfigYamlFile = "env.yaml"
	DevConfigYamlFile = "dev.yaml"
	ConfigYamlFile    = "config.yaml"
)

// Build and deploy
var (
	//go:embed pattern_c/Dockerfile
	dockerfile []byte
	Dockerfile = &folder.Folder{
		Name:    "Dockerfile",
		Content: dockerfile,
	}

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
	//go:embed pattern_c/Makefile
	makefile []byte
	Makefile = &folder.Folder{
		Name:    "Makefile",
		Content: makefile,
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

// Testing files
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
	//go:embed pattern_c/scripts/rscli.mk
	RscliMK []byte

	//go:embed pattern_c/scripts/grpc.mk
	GrpcMK []byte
)

var (
	//go:embed pattern_c/api/grpc/api.proto
	protoServer []byte
	ProtoServer = &folder.Folder{
		Name:    "api.proto",
		Content: protoServer,
	}
)
