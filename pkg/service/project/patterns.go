package project

import (
	_ "embed"
)

var datasourceClients = map[dataSourcePrefix][]byte{}

func init() {
	datasourceClients[RedisDataSourcePrefix] = redisConn
	datasourceClients[PostgresDataSourcePrefix] = pgConn
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

//go:embed pattern_c/docker-compose.yml
var dockerCompose []byte

//go:embed pattern_c/README.md
var readme []byte
