package patterns

import (
	"bytes"
	_ "embed"

	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/pkg/folder"
	_const "github.com/Red-Sock/rscli/plugins/config/pkg/const"
)

// Constants naming: Purpose+Type (File)

const (
	ImportProjectNamePatternKebabCase = "financial-microservice"
	ImportProjectNamePatternSnakeCase = "financial_microservice"
)

var DatasourceClients = map[string][]*folder.Folder{}
var ServerOptsPatterns = map[string]serverPattern{}

type serverPattern struct {
	F          folder.Folder
	Validators func(f *folder.Folder, serverName string)
}

const (
	ServerGoFile  = "server.go"
	versionGoFile = "version.go"
	pingerGoFile  = "pinger.go"

	HandlerFolder = "handlers"
	HandlerGoFile = "handler.go"
)

func init() {
	DatasourceClients[_const.SourceNameRedis] = []*folder.Folder{{Name: ConnFile, Content: RedisConnFile}}
	DatasourceClients[_const.SourceNamePostgres] = []*folder.Folder{{Name: ConnFile, Content: PgConnFile}, {Name: PgTxFileName, Content: PgTxFile}}
	DatasourceClients[_const.TelegramServer] = []*folder.Folder{{Name: ConnFile, Content: TgConnFile}}

	ServerOptsPatterns[_const.RESTHTTPServer] = serverPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				{
					Name:    ServerGoFile,
					Content: RestServFile,
				},
				{
					Name:    versionGoFile,
					Content: RestServHandlerExampleFile,
				},
			},
		},
		Validators: func(f *folder.Folder, serverName string) {
			for _, innerFolder := range f.Inner {
				innerFolder.Content = bytes.ReplaceAll(
					innerFolder.Content,
					[]byte("package rest_realisation"),
					[]byte("package "+serverName),
				)

				if innerFolder.Name == ServerGoFile {
					innerFolder.Content = bytes.ReplaceAll(
						innerFolder.Content,
						[]byte("config.ServerRestApiPort"),
						[]byte("config.Server"+cases.SnakeToCamel(serverName)+"Port"))
				}
			}
		},
	}

	ServerOptsPatterns[_const.TelegramServer] = serverPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				{
					Name:    ServerGoFile,
					Content: TgServFile,
				},
				{
					Name: HandlerFolder,
					Inner: []*folder.Folder{
						{
							Name: "version",
							Inner: []*folder.Folder{
								{
									Name:    HandlerGoFile,
									Content: TgHandlerExampleFile,
								},
							},
						},
					},
				},
			},
		},
		Validators: func(f *folder.Folder, serverName string) {
			server := f.GetByPath(ServerGoFile)

			server.Content = bytes.ReplaceAll(
				server.Content,
				[]byte("package tg"),
				[]byte("package "+serverName),
			)
			server.Content = bytes.ReplaceAll(
				server.Content,
				[]byte("config.ServerTgApikey"),
				[]byte("config.Server"+cases.SnakeToCamel(serverName)+"Apikey"))
		},
	}

	ServerOptsPatterns[_const.GRPCServer] = serverPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				{
					Name:    ServerGoFile,
					Content: GrpcServFile,
				},
				{
					Name:    pingerGoFile,
					Content: GrpcServExampleFile,
				},
			},
		},
		Validators: func(f *folder.Folder, serverName string) {
			for _, innerFolder := range f.Inner {
				innerFolder.Content = bytes.ReplaceAll(
					innerFolder.Content,
					[]byte("package grpc_realisation"),
					[]byte("package "+serverName),
				)

				innerFolder.Content = bytes.ReplaceAll(
					innerFolder.Content,
					[]byte("/pkg/grpc-realisation\""),
					[]byte("/pkg/"+serverName+"\""),
				)

				if innerFolder.Name == ServerGoFile {
					innerFolder.Content = bytes.ReplaceAll(
						innerFolder.Content,
						[]byte("config.ServerGRPCApiPort"),
						[]byte("config.Server"+cases.SnakeToCamel(serverName)+"Port"))
				}
			}
		},
	}
}

const (
	CmdFolder    = "cmd"
	MainFileName = "main.go"

	BootStrapFolder = "bootstrap"

	ApiConstructorFileName = "api.go"

	InternalFolder = "internal"
	ClientsFolder  = "clients"
	PostgresFolder = "postgres"
	ConnFile       = "conn.go"
	PgTxFileName   = "tx.go"

	PkgFolder          = "pkg"
	ProtoFolder        = "proto"
	ProtoFileExtension = ".proto"

	UtilsFolder  = "utils"
	CloserFolder = "closer"
	CloserFile   = "closer.go"

	TransportFolder    = "transport"
	ApiManagerFileName = "manager.go"

	ConfigsFolder  = "config"
	ConfigTemplate = "config.yaml.template"

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
