package consts

type DataSourcePrefix string

const (
	RedisDataSourcePrefix    DataSourcePrefix = "redis"
	PostgresDataSourcePrefix DataSourcePrefix = "postgres"
)

type ServerOptsPrefix string

const (
	RESTServerPrefix ServerOptsPrefix = "rest"
	GRPCServerPrefix ServerOptsPrefix = "grpc"
)
