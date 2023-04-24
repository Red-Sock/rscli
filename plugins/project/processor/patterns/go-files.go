package patterns

import (
	"bytes"
	_ "embed"

	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/pkg/folder"
	_const "github.com/Red-Sock/rscli/plugins/config/pkg/const"
)

var DatasourceClients = map[string][]byte{}
var ServerOptsPatterns = map[string]serverPattern{}

type serverPattern struct {
	F          folder.Folder
	Validators func(f *folder.Folder, serverName string)
}

const (
	ServerGoFile  = "server.go"
	versionGoFile = "version.go"

	MenuFolder    = "menus"
	MenuGoFile    = "menu.go"
	HandlerFolder = "handlers"
	HandlerGoFile = "handler.go"
)

func init() {
	DatasourceClients[_const.SourceNameRedis] = RedisConn
	DatasourceClients[_const.SourceNamePostgres] = PgConn

	ServerOptsPatterns[_const.RESTHTTPServer] = serverPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				{
					Name:    ServerGoFile,
					Content: RestServ,
				},
				{
					Name:    versionGoFile,
					Content: RestServVersion,
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
					Content: TgServ,
				},
				{
					Name: MenuFolder,
					Inner: []*folder.Folder{
						{
							Name: "mainmenu",
							Inner: []*folder.Folder{
								{
									Name:    MenuGoFile,
									Content: TgMainMenu,
								},
							},
						},
					},
				},
				{
					Name: HandlerFolder,
					Inner: []*folder.Folder{
						{
							Name: "version",
							Inner: []*folder.Folder{
								{
									Name:    HandlerGoFile,
									Content: TgVersionHandler,
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

	ServerOptsPatterns[_const.GRPCServer] = serverPattern{}
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

	UtilsFolder  = "utils"
	CloserFolder = "closer"
	CloserFile   = "closer.go"

	TransportFolder    = "transport"
	ApiManagerFileName = "manager.go"

	ConfigsFolder  = "config"
	ConfigTemplate = "config.yaml.template"
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
	RedisConn []byte
	//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
	PgConn []byte
)

// Config parser files
var (
	//go:embed pattern_c/internal/config/config.go.pattern
	Configurator string
	//go:embed pattern_c/internal/config/keys.go.pattern
	ConfigKeys []byte
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	ServerManagerPattern []byte

	//go:embed pattern_c/internal/transport/rest_realisation/server.go.pattern
	RestServ []byte
	//go:embed pattern_c/internal/transport/rest_realisation/version.go.pattern
	RestServVersion []byte

	//go:embed pattern_c/internal/transport/tg/listener.go.pattern
	TgServ []byte
	//go:embed pattern_c/internal/transport/tg/menus/mainmenu/main-menu.go.pattern
	TgMainMenu []byte
	//go:embed pattern_c/internal/transport/tg/handlers/version/handler.go.pattern
	TgVersionHandler []byte

	// TODO
	GrpcServ []byte
)

// Utils
var (
	//go:embed pattern_c/internal/utils/closer/closer.go.pattern
	UtilsCloser []byte
)
