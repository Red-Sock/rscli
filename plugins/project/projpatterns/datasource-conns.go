package projpatterns

import (
	"bytes"
	"fmt"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/cases"
)

func GetDatasourceClientFile(datasourceName string) (*folder.Folder, error) {
	out := &folder.Folder{Name: datasourceName}
	switch datasourceName {
	case SourceNameRedis:
		out.Inner = []*folder.Folder{RedisConnFile.Copy()}
	case SourceNamePostgres:
		out.Inner = []*folder.Folder{PgConnFile.Copy(), PgTxFile.Copy()}
	case TelegramServer:
		out.Inner = []*folder.Folder{TgConnFile.Copy()}
	default:
		return nil, errors.New(fmt.Sprintf("unknown data source %s. "+
			"DataSource should start with name of source (e.g redis, postgres)"+
			"and (or) be followed by \"_\" + actual_unique_name symbol if needed "+
			"(e.g redis_shard1, postgres_replica2)", datasourceName))
	}

	return out, nil
}

type ServerPattern struct {
	F          folder.Folder
	Validators func(f *folder.Folder, serverName string)
}

func GetServerFiles(serverName string) (ServerPattern, error) {
	switch serverName {
	case RESTHTTPServer:
		return getRestPattern(), nil
	case GRPCServer:
		return getGrpcPattern(), nil
	case TelegramServer:
		return getTelegramPattern(), nil

	default:
		return ServerPattern{}, errors.New(fmt.Sprintf("unknown server option %s. ", serverName) +
			"Server Option should start with type of server (e.g rest, grpc)" +
			"and (or) be followed by \"_\" symbol + unique_name if needed (e.g rest_v1, grpc_proxy)")
	}

}

func getRestPattern() ServerPattern {
	return ServerPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				RestServFile.Copy(),
				RestServHandlerVersionExampleFile.Copy(),
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
}

func getTelegramPattern() ServerPattern {
	return ServerPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				TgServFile.Copy(),
				{
					Name: handlerFolder,
					Inner: []*folder.Folder{
						{
							Name: "version",
							Inner: []*folder.Folder{
								TgHandlerExampleFile.Copy(),
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
				[]byte("config.ServerTelegramAPIKey"),
				[]byte("config.Server"+cases.SnakeToCamel(serverName)+"Apikey"))
		},
	}
}

func getGrpcPattern() ServerPattern {
	return ServerPattern{
		F: folder.Folder{
			Inner: []*folder.Folder{
				GrpcServFile.Copy(),
				GrpcServExampleFile.Copy(),
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
