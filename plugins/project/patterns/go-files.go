package patterns

import (
	_ "embed"
)

// Constants naming: Purpose+Type (File)

const (
	ImportProjectNamePatternKebabCase = "financial-microservice"
	ImportProjectNamePatternSnakeCase = "financial_microservice"
)

const (
	serverGoFile  = "server.go"
	versionGoFile = "version.go"
	pingerGoFile  = "pinger.go"

	handlerFolder = "handlers"
	handlerGoFile = "handler.go"
)

const (
	CmdFolder    = "cmd"
	MainFileName = "main.go"

	BootStrapFolder = "bootstrap"

	ApiFolder              = "api"
	ApiConstructorFileName = "api.go"

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

	ConfigsFolder      = "config"
	ConfigFileName     = "config.go"
	ConfigTemplate     = "config.yaml.template"
	ConfigKeysFileName = "keys.go"

	GoMod = "go.mod"
)

// Basic files
var (
	//go:embed pattern_c/cmd/financial-microservice/main.go.pattern
	MainFile []byte
	//go:embed pattern_c/cmd/financial-microservice/bootstrap/api.go.pattern
	APISetupFile []byte
)

// DataStorage connection files
var (
	//go:embed pattern_c/internal/clients/redis/conn.go.pattern
	RedisConnFile []byte
	//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
	PgConnFile []byte
	//go:embed pattern_c/internal/clients/postgres/tx.go.pattern
	PgTxFile []byte
	//go:embed pattern_c/internal/clients/telegram/conn.go.pattern
	TgConnFile []byte
)

// Config parser files
var (
	//go:embed pattern_c/internal/config/config.go.pattern
	ConfiguratorFile string
	//go:embed pattern_c/internal/config/keys.go.pattern
	ConfigKeysFile []byte
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	ServerManagerPatternFile []byte

	//go:embed pattern_c/internal/transport/rest_realisation/server.go.pattern
	RestServFile []byte
	//go:embed pattern_c/internal/transport/rest_realisation/version.go.pattern
	RestServHandlerExampleFile []byte

	//go:embed pattern_c/internal/transport/tg/listener.go.pattern
	TgServFile []byte
	//go:embed pattern_c/internal/transport/tg/handlers/version/handler.go.pattern
	TgHandlerExampleFile []byte

	//go:embed pattern_c/internal/transport/grpc_realisation/server.go.pattern
	GrpcServFile []byte
	//go:embed pattern_c/internal/transport/grpc_realisation/pinger.go.pattern
	GrpcServExampleFile []byte
	//go:embed pattern_c/pkg/proto/grpc_realisation/financial-microservice.proto
	GrpcProtoExampleFile []byte
)

// Utils
var (
	//go:embed pattern_c/internal/utils/closer/closer.go.pattern
	UtilsCloserFile []byte
)
