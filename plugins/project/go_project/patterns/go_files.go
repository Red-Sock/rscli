package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Constants naming: Purpose+Type (File)

const (
	GithubFolder    = ".github"
	WorkflowsFolder = "workflows"

	CmdFolder     = "cmd"
	ServiceFolder = "service"

	InternalFolder             = "internal"
	AppFolder                  = "app"
	AppFileName                = "app.go"
	AppInitServerFileName      = "server.go"
	AppInitDataSourcesFileName = "data_sources.go"
	AppConfigFileName          = "config.go"
	AppCustomFileName          = "custom.go"

	ConnFileName = "conn.go"

	TransportFolder = "transport"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder      = "config"
	ConfigTemplateYaml = "config_template.yaml"

	ConfigLoadFileName        = "load.go"
	ConfigDataSourcesFileName = "data_sources.go"
	ConfigEnvironmentFileName = "environment.go"

	GoMod = "go.mod"

	ExampleFile = ".example"
)

// Basic files
var (
	//go:embed pattern_c/cmd/service/main.go.pattern
	mainFile []byte
	MainFile = &folder.Folder{
		Name:    "main.go",
		Content: mainFile,
	}
)

// Clients connection files
var (
	//go:embed pattern_c/internal/clients/redis/conn.go.pattern
	redisConnFile []byte
	RedisConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: redisConnFile,
	}

	//go:embed pattern_c/internal/clients/sqldb/conn.go.pattern
	sqlConnFile []byte
	SqlConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: sqlConnFile,
	}

	//go:embed pattern_c/internal/clients/telegram/conn.go.pattern
	tgConnFile []byte
	TgConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: tgConnFile,
	}

	//go:embed pattern_c/internal/clients/grpc/conn.go.pattern
	grpcClientConnFile []byte
	GrpcClientConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: grpcClientConnFile,
	}
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	serverManagerPatternFile []byte
	ServerManager            = &folder.Folder{
		Name:    "manager.go",
		Content: serverManagerPatternFile,
	}

	//go:embed pattern_c/internal/transport/grpc.go.pattern
	grpcServerManagerPatternFile []byte
	GrpcServerManagerPatternFile = &folder.Folder{
		Name:    "grpc.go",
		Content: grpcServerManagerPatternFile,
	}

	//go:embed pattern_c/internal/transport/http.go.pattern
	httpServerManagerPatternFile []byte
	HttpServerManagerPatternFile = &folder.Folder{
		Name:    "http.go",
		Content: httpServerManagerPatternFile,
	}

	//go:embed pattern_c/internal/transport/telegram/listener.go.pattern
	tgServFile []byte
	TgServFile = &folder.Folder{
		Name:    "listener.go",
		Content: tgServFile,
	}
	//go:embed pattern_c/internal/transport/telegram/version/handler.go.pattern
	tgHandlerExampleFile []byte
	TgHandlerExampleFile = &folder.Folder{
		Name:    "handler.go",
		Content: tgHandlerExampleFile,
	}
)
