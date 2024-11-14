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
