package config

type configKey string

const (
	// _start_of_consts_to_replace

	AppInfoVersion             = "app_info_version"
	AppInfoStartupDuration     = "app_info_startup_deadline"
	ServerRestApiPort          = "server_rest_api_port"
	ServerTgApikey             = "server_tg_api_key"
	DataSourcesPostgresUser    = "data_sources_postgres_user"
	DataSourcesPostgresPwd     = "data_sources_postgres_pwd"
	DataSourcesPostgresHost    = "data_sources_postgres_host"
	DataSourcesPostgresPort    = "data_sources_postgres_port"
	DataSourcesPostgresName    = "data_sources_postgres_name"
	DataSourcesPostgresSslmode = "data_sources_postgres_connection_sslmode"
	// _end_of_consts_to_replace
)
