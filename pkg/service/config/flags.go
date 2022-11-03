package config

const (
	DataSourceKey = "data_sources"
	ServerOptsKey = "server"
)

// data sources (sub-keys)
const (
	SourceNamePg  = "pg"
	SourceNameRds = "rds"
)

const (
	RESTHTTPServer = "rest"
	GRPCServer     = "grpc"
)

// flags
const (
	forceOverride = "fo"
	configPath    = "path"
)

const (
	DefaultDir = "config"
	FileName   = "config.yaml"
)
