package shared

import (
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/pkg/flag/flags"
	"os"
	"path"
)

const pluginDir = "rscli_plugins"

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

	pluginsDir, err := os.Executable()
	if err != nil {
		panic(err)
	}

	pluginsDir, _ = path.Split(pluginsDir)

	return path.Join(pluginsDir, pluginDir)
}
