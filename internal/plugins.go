package internal

import (
	cfgcmd "github.com/Red-Sock/rscli/plugins/cfg/cmd"
	cfgui "github.com/Red-Sock/rscli/plugins/cfg/ui"
)

var (
	plugins = map[string]Plugin{
		cfgcmd.PluginName: &cfgcmd.Plug,
	}
	pluginsWithUI = map[string]PluginWithUi{
		cfgui.PluginName: &cfgui.Plug,
	}
)
