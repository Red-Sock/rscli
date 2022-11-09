package config

// names of core items in final config
const (
	DataSourceKey = "data_sources"
	ServerOptsKey = "server"
	AppKey        = "app"
)

// data sources (sub-keys) flags
const (
	SourceNamePg  = "pg"
	SourceNameRds = "rds"
)

// server type flags
const (
	RESTHTTPServer = "rest"
	GRPCServer     = "grpc"
)

const (
	AppName    = "app_name"
	AppVersion = "app_version"
)

// additional flags
const (
	forceOverride = "fo"
	configPath    = "path"
)

// default values
const (
	DefaultDir = "config"
	FileName   = "config.yaml"
)
