package _const

// names of core items in final config
const (
	DataSourceKey = "data_sources"
	ServerOptsKey = "server"
	AppKey        = "app"
)

// data sources (sub-keys) flags
const (
	SourceNamePostgres = "postgres"
	SourceNameRedis    = "redis"
)

// server type flags
const (
	RESTHTTPServer = "rest"
	GRPCServer     = "grpc"
	TelegramServer = "tg"
)

const (
	AppName    = "app_name"
	AppVersion = "app_version"
)

// additional flags
const (
	ForceOverride = "fo"
	ConfigPath    = "path"
)

// default values
const (
	DefaultDir = "config"
	FileName   = "dev.yaml"
)
