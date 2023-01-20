package internal

import (
	cfgcmd "github.com/Red-Sock/rscli/plugins/config/cmd"
	cfgui "github.com/Red-Sock/rscli/plugins/config/ui"

	projectcmd "github.com/Red-Sock/rscli/plugins/project/cmd"
	projectui "github.com/Red-Sock/rscli/plugins/project/ui"
)

var (
	plugins = map[string]Plugin{
		cfgcmd.PluginName:     &cfgcmd.Plug,
		projectcmd.PluginName: &projectcmd.Plug,
	}
	pluginsWithUI = map[string]PluginWithUi{
		cfgui.PluginName:     &cfgui.Plug,
		projectui.PluginName: &projectui.Plug,
	}
)
