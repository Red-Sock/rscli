package add

import (
	_ "embed"
)

// GRPC
var (
	//go:embed expected/grpc/proto.proto
	grpcExpectedProtoFile []byte
	//go:embed expected/grpc/config.yaml
	grpcMatreshkaConfigExpected []byte
)

// Redis
var (
	//go:embed expected/redis/config.yaml
	expectedRedisConfig []byte

	//go:embed expected/redis/config/data_source.go
	expectedRedisDataSourceConfig []byte
	//go:embed expected/redis/app/data_source.go
	expectedRedisDataSourceApp []byte
)

// Postgres
var (
	//go:embed expected/postgres/config.yaml
	expectedPostgresConfig []byte
	//go:embed expected/postgres/config/data_sources.go
	expectedPostgresDataSourceConfig []byte
	//go:embed expected/postgres/app/data_sources.go
	expectedPostgresDataSourceApp []byte
)

// Telegram
var (
	//go:embed expected/telegram/config.yaml
	expectedTelegramConfig []byte
	//go:embed expected/telegram/app/data_sources.go
	expectedTelegramDataSourcesApp []byte
	//go:embed expected/telegram/config/data_sources.go
	expectedTelegramDataSourcesConfig []byte
	//go:embed expected/telegram/transport/listener.go
	expectedTelegramServer []byte
	//go:embed expected/telegram/transport/handler/handler.go
	expectedTelegramServerHandlerExample []byte
)
