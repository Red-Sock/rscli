package shared

import (
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/pkg/flag/flags"
	"os"
)

func GetPluginsDir(flgs map[string][]string) string {
	if flgs != nil && len(flgs) != 0 {
		if pluginsDir, err := flag.ExtractOneValueFromFlags(flgs, flags.PluginsDirFlag); err != nil && pluginsDir != "" {
			return pluginsDir
		}
	}

	pluginsDir, ok := os.LookupEnv(flags.PluginsDirEnv)
	if ok {
		return pluginsDir
	}

	return ""
}
