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

//go:embed pattern_c/internal/transport/rest/server.go.pattern
var restServ []byte

//go:embed pattern_c/internal/transport/rest/version.go.pattern
var restServVersion []byte

// TODO
var grpcServ []byte
