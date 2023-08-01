package processor

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/plugins/project/config/pkg/configstructs"
	"github.com/Red-Sock/rscli/plugins/project/config/pkg/const"
)

var ErrPatternExists = errors.New("pattern exists")

func parseFlag(f string, args []string, cfg *configstructs.Config) error {

	if strings.HasPrefix(f, "-") {
		f = f[1:]
	}

	switch f {
	// data sources
	case _const.SourceNamePostgres:
		return addPattern(DefaultPgPattern(args), cfg.DataSources)
	case _const.SourceNameRedis:
		return addPattern(DefaultRdsPattern(args), cfg.DataSources)

	// transport layer
	case _const.RESTHTTPServer:
		return addPattern(DefaultHTTPPattern(args), cfg.Server)
	case _const.GRPCServer:
		return addPattern(DefaultGRPCPattern(args), cfg.Server)
	case _const.TelegramServer:
		return addPattern(DefaultTelegramPattern(args), cfg.Server)

	// app info
	case _const.AppName:
		if len(args) != 1 {
			return errors.New("INVALID ARGUMENTS AMOUNT FOR FLAG " + _const.AppName)
		}
		cfg.AppInfo.Name = args[0]
		return nil
	case _const.AppVersion:
		if len(args) != 1 {
			return errors.New("INVALID ARGUMENTS AMOUNT FOR FLAG " + _const.AppVersion)
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
			return errors.Wrapf(ErrPatternExists, "pattern with name %s exists", k)
		}
	}

	for k, v := range src {
		tgt[k] = v
	}

	return nil
}
