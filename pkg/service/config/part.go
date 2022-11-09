package config

import "strings"

func parseFlag(f string, args []string) (name string, values map[string]interface{}) {

	if strings.HasPrefix(f, "-") {
		f = f[1:]
	}

	switch f {
	// data sources
	case SourceNamePg:
		return DataSourceKey, DefaultPgPattern(args)
	case SourceNameRds:
		return DataSourceKey, DefaultRdsPattern(args)

	// transport layer
	case RESTHTTPServer:
		return ServerOptsKey, DefaultHTTPPattern(args)
	case GRPCServer:
		return ServerOptsKey, DefaultGRPCPattern(args)

	// app info
	case AppName:
		return AppKey, AppNamePattern(args)
	case AppVersion:
		return AppKey, AppVersionPattern(args)
	default:
		return "", nil
	}
}
