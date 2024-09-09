package projpatterns

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

	ConfigsFolder  = "config"
	ConfigTemplate = "config_template.yaml"

	ConfigLoadFileName        = "load.go"
	ConfigDataSourcesFileName = "data_sources.go"
	ConfigEnvironmentFileName = "environment.go"

	GoMod = "go.mod"

	ExampleFile = ".example"
)

// Basic files
var (
	//go:embed pattern_c/cmd/rscli_example/main.go.pattern
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

// Config parser files
var (
	//go:embed pattern_c/internal/config/autoload.go.pattern
	autoloadConfigFile []byte
	AutoloadConfigFile = &folder.Folder{
		Name:    "autoload.go",
		Content: autoloadConfigFile,
	}
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	serverManagerPatternFile []byte
	ServerManagerPatternFile = &folder.Folder{
		Name:    "manager.go",
		Content: serverManagerPatternFile,
	}

	//go:embed pattern_c/internal/transport/rest/server.go.pattern
	restServFile []byte
	RestServFile = &folder.Folder{
		Name:    "server.go",
		Content: restServFile,
	}
	//go:embed pattern_c/internal/transport/rest/version.go.pattern
	restServHandlerVersionExampleFile []byte
	RestServHandlerVersionExampleFile = &folder.Folder{
		Name:    "version.go",
		Content: restServHandlerVersionExampleFile,
	}

	//go:embed pattern_c/internal/transport/grpc/server.go.pattern
	grpcServFile []byte
	GrpcServFile = &folder.Folder{
		Name:    "server.go",
		Content: grpcServFile,
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
