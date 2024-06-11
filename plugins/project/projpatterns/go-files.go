package projpatterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Constants naming: Purpose+Type (File)

const (
	CmdFolder     = "cmd"
	ServiceFolder = "service"

	ExamplesFolder = "examples"

	ExamplesHttpEnvFile = "http-client.env.json"

	InternalFolder = "internal"
	ConnFileName   = "conn.go"
	PgTxFileName   = "tx.go"

	TransportFolder = "transport"

	PkgFolder = "pkg"

	UtilsFolder  = "utils"
	CloserFolder = "closer"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder             = "config"
	ConfigTemplate            = "config_template.yaml"
	ConfigEnvironmentFileName = "environment.go"
	ConfigKeysFileName        = "keys.go"

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

	//go:embed pattern_c/internal/clients/grpc/conn.go.pattern
	grpcClientConnFile []byte
	GrpcClientConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: grpcClientConnFile,
	}

	//go:embed pattern_c/internal/clients/sqlite/conn.go.pattern
	sqliteClientConnFile []byte
	SqliteClientConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: sqliteClientConnFile,
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

// Utils
var (
	//go:embed pattern_c/internal/utils/closer/closer.go.pattern
	utilsCloserFile []byte
	UtilsCloserFile = &folder.Folder{
		Name:    "closer.go",
		Content: utilsCloserFile,
	}
)
