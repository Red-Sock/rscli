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
	ServerGoFile  = "server.go"
	VersionGoFile = "version.go"
	PingerGoFile  = "pinger.go"

	handlerFolder = "handlers"
	handlerGoFile = "handler.go"
)

const (
	CmdFolder = "cmd"

	BootStrapFolder = "bootstrap"

	ApiFolder = "api"

	ExamplesFolder      = "examples"
	ExampleFileName     = "api.http"
	ExamplesHttpEnvFile = "http-client.env.json"

	InternalFolder = "internal"
	ClientsFolder  = "clients"
	PostgresFolder = "postgres"
	ConnFileName   = "conn.go"
	PgTxFileName   = "tx.go"

	PkgFolder          = "pkg"
	SwaggerFolder      = "swagger"
	ProtoFolder        = "proto"
	ProtoFileExtension = ".proto"

	UtilsFolder  = "utils"
	CloserFolder = "closer"
	CloserFile   = "closer.go"

	TransportFolder    = "transport"
	ApiManagerFileName = "manager.go"

	HandlersFolderName   = "handlers"
	VersionFolderName    = "version"
	TelegramServFileName = "listener.go"
	TgHandlerFileName    = "handler.go"
	ConfigsFolder        = "config"
	ConfigFileName       = "config.go"
	ConfigTemplate       = "config_template.yaml"

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
		Name:    "conn.go",
		Content: redisConnFile,
	}

	//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
	pgConnFile []byte
	PgConnFile = &folder.Folder{
		Name:    "conn.go",
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
		Name:    "conn.go",
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

	//go:embed pattern_c/internal/config/keys.go.pattern
	configKeysFile []byte
	ConfigKeysFile = &folder.Folder{
		Name:    "keys.go",
		Content: configKeysFile,
	}
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	ServerManagerPatternFile []byte

	//go:embed pattern_c/internal/transport/rest_api/server.go.pattern
	RestServFile []byte
	//go:embed pattern_c/internal/transport/rest_api/version.go.pattern
	RestServHandlerExampleFile []byte

	//go:embed pattern_c/internal/transport/tg/listener.go.pattern
	TgServFile []byte
	//go:embed pattern_c/internal/transport/tg/handlers/version/handler.go.pattern
	TgHandlerExampleFile []byte

	//go:embed pattern_c/internal/transport/grpc_api/server.go.pattern
	GrpcServFile []byte
	//go:embed pattern_c/internal/transport/grpc_api/pinger.go.pattern
	GrpcServExampleFile []byte
	//go:embed pattern_c/pkg/proto/grpc_realisation/financial-microservice.proto
	GrpcProtoExampleFile []byte
)

// Utils
var (
	//go:embed pattern_c/internal/utils/closer/closer.go.pattern
	UtilsCloserFile []byte
)
