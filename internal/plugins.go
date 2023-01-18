package internal

import (
	"github.com/Red-Sock/rscli/internal/utils/embeded"
)

var (
	plugins       = map[string]Plugin{}
	pluginsWithUI = map[string]PluginWithUi{}

	basicPlugin = []Plugin{
		&embeded.GetPlugin{},
		&embeded.DeletePlugin{},
	}
)
