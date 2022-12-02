package project

import (
	_ "embed"
)

var datasourceClients = map[dataSourcePrefix][]byte{}
var serverOptsPatterns = map[serverOptsPrefix]map[string][]byte{}

func init() {
	datasourceClients[RedisDataSourcePrefix] = redisConn
	datasourceClients[PostgresDataSourcePrefix] = pgConn

	serverOptsPatterns[RESTServerPrefix] = map[string][]byte{"server.go": restServ, "version.go": restServVersion}
	// TODO
	serverOptsPatterns[GRPCServerPrefix] = map[string][]byte{}
}

//go:embed pattern_c/cmd/main.go.pattern
var mainFile []byte

//go:embed pattern_c/internal/clients/redis/conn.go.pattern
var redisConn []byte

//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
var pgConn []byte

//go:embed pattern_c/internal/config/config.go.pattern
var configurator string

//go:embed pattern_c/Dockerfile
var dockerfile []byte

//go:embed pattern_c/README.md
var readme []byte

//go:embed pattern_c/internal/transport/rest_realisation/server.go.pattern
var restServ []byte

//go:embed pattern_c/internal/transport/rest_realisation/version.go.pattern
var restServVersion []byte

// TODO
var grpcServ []byte

//go:embed pattern_c/internal/transport/manager.go.pattern
var managerPattern []byte

//go:embed pattern_c/examples/api.http
var apiHTTP []byte

//go:embed pattern_c/examples/http-client.env.json
var httpEnvironment []byte
