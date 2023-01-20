package processor

import (
	"errors"
	"github.com/Red-Sock/rscli/plugins/config/pkg/structs"
	"strings"
)

func parseFlag(f string, args []string, cfg *structs.Config) error {

	if strings.HasPrefix(f, "-") {
		f = f[1:]
	}

	switch f {
	// data sources
	case SourceNamePg:
		return addPattern(DefaultPgPattern(args), cfg.DataSources)
	case SourceNameRds:
		return addPattern(DefaultRdsPattern(args), cfg.DataSources)

	// transport layer
	case RESTHTTPServer:
		return addPattern(DefaultHTTPPattern(args), cfg.Server)
	case GRPCServer:
		return addPattern(DefaultGRPCPattern(args), cfg.Server)

	// app info
	case AppName:
		if len(args) != 1 {
			return errors.New("INVALID ARGUMENTS AMOUNT FOR FLAG " + AppName)
		}
		cfg.AppInfo.Name = args[0]
		return nil
	case AppVersion:
		if len(args) != 1 {
			return errors.New("INVALID ARGUMENTS AMOUNT FOR FLAG " + AppVersion)
		}
		cfg.AppInfo.Version = args[0]
		return nil
	default:
		return errors.New("UNKNOWN FLAG " + f)
	}
}

func addPattern(src, tgt map[string]interface{}) error {
	for k := range src {
		if _, ok := tgt[k]; ok {
			return errors.New("")
		}
	}

	for k, v := range src {
		tgt[k] = v
	}

	return nil
}
