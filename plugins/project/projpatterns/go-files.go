package projpatterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Constants naming: Purpose+Type (File)

const (
	ImportProjectNamePatternKebabCase = "financial-microservice"
	ImportProjectNamePatternSnakeCase = "financial_microservice"
)

const (
	ServerGoFile = "server.go"

	handlerFolder = "handlers"
)

const (
	CmdFolder = "cmd"

	BootStrapFolder = "bootstrap"

	ApiFolder = "api"

	ExamplesFolder = "examples"

	ExamplesHttpEnvFile = "http-client.env.json"

	InternalFolder = "internal"
	ConnFileName   = "conn.go"
	PgTxFileName   = "tx.go"

	PkgFolder     = "pkg"
	SwaggerFolder = "swagger"

	UtilsFolder  = "utils"
	CloserFolder = "closer"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder      = "config"
	ConfigTemplate     = "config_template.yaml"
	ConfigKeysFileName = "keys.go"

	GoMod = "go.mod"

	ExampleFile = ".example"
)

// Basic files
var (
	//go:embed pattern_c/cmd/financial-microservice/main.go.pattern
	mainFile []byte
	MainFile = &folder.Folder{
		Name:    "main.go",
		Content: mainFile,
	}

	//go:embed pattern_c/cmd/financial-microservice/bootstrap/api.go.pattern
	apiSetupFile []byte
	APISetupFile = &folder.Folder{
		Name:    "api.go",
		Content: apiSetupFile,
	}
)

// DataStorage connection files
var (
	//go:embed pattern_c/internal/clients/redis/conn.go.pattern
	redisConnFile []byte
	RedisConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: redisConnFile,
	}

	//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
	pgConnFile []byte
	PgConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: pgConnFile,
	}
	//go:embed pattern_c/internal/clients/postgres/tx.go.pattern
	pgTxFile []byte
	PgTxFile = &folder.Folder{
		Name:    "tx.go",
		Content: pgTxFile,
	}

	//go:embed pattern_c/internal/clients/telegram/conn.go.pattern
	tgConnFile []byte
	TgConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: tgConnFile,
	}
)

// Config parser files
var (
	//go:embed pattern_c/internal/config/config.go.pattern
	configFile []byte
	ConfigFile = &folder.Folder{
		Name:    "config.go",
		Content: configFile,
	}
	//go:embed pattern_c/internal/config/autoload.go.pattern
	autoloadConfigFile []byte
	AutoloadConfigFile = &folder.Folder{
		Name:    "autoload.go",
		Content: autoloadConfigFile,
	}
	//go:embed pattern_c/internal/config/static.go.pattern
	staticConfigFile []byte
	StaticConfigFile = &folder.Folder{
		Name:    "static.go",
		Content: staticConfigFile,
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

	//go:embed pattern_c/internal/transport/rest_api/server.go.pattern
	restServFile []byte
	RestServFile = &folder.Folder{
		Name:    ServerGoFile,
		Content: restServFile,
	}
	//go:embed pattern_c/internal/transport/rest_api/version.go.pattern
	restServHandlerVersionExampleFile []byte
	RestServHandlerVersionExampleFile = &folder.Folder{
		Name:    "version.go",
		Content: restServHandlerVersionExampleFile,
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

	//go:embed pattern_c/internal/transport/grpc_api/server.go.pattern
	grpcServFile []byte
	GrpcServFile = &folder.Folder{
		Name:    ServerGoFile,
		Content: grpcServFile,
	}
	//go:embed pattern_c/internal/transport/grpc_api/pinger.go.pattern
	grpcServExampleFile []byte
	GrpcServExampleFile = &folder.Folder{
		Name:    "pinger.go",
		Content: grpcServExampleFile,
	}
	//go:embed pattern_c/pkg/proto/grpc_realisation/financial-microservice.proto
	grpcProtoExampleFile []byte
	GrpcProtoExampleFile = &folder.Folder{
		Name:    "financial-microservice.proto",
		Content: grpcProtoExampleFile,
	}
)

// Utils
var (
	//go:embed pattern_c/internal/utils/closer/closer.go.pattern
	utilsCloserFile []byte
	UtilsCloserFile = &folder.Folder{
		Name:    "closer.go",
		Content: utilsCloserFile,
	}
)
