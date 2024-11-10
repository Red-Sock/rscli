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

	GenCommand           = "gen"
	GenGrpcServerCommand = "gen-server-grpc"
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
	//go:embed pattern_c/rscli.mk
	rscliMK []byte
	RscliMK = &folder.Folder{
		Name:    "rscli.mk",
		Content: rscliMK,
	}
	//go:embed pattern_c/grpc.mk
	GrpcServerGenMK []byte
)

var (
	//go:embed pattern_c/api/grpc/api.proto
	protoServer []byte
	ProtoServer = &folder.Folder{
		Name:    "api.proto",
		Content: protoServer,
	}
)

// GitHub Workflows
var (
	//go:embed pattern_c/.github/workflows/release.yaml
	githubWorkflowRelease []byte
	GithubWorkflowRelease = &folder.Folder{
		Name:    "release.yaml",
		Content: githubWorkflowRelease,
	}

	//go:embed pattern_c/.github/workflows/go-branch-push.yml
	githubWorkflowGoBranchPush []byte
	GithubWorkflowGoBranchPush = &folder.Folder{
		Name:    "branch-push.yaml",
		Content: githubWorkflowGoBranchPush,
	}
)

// Scripts
var (
	//go:embed pattern_c/Makefile
	makefile []byte
	Makefile = &folder.Folder{
		Name:    MakefileFile,
		Content: makefile,
	}
)
