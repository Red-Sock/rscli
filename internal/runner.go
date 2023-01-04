package internal

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/flag"

	"github.com/Red-Sock/rscli/pkg/service/help"
)

const (
	pluginsDIR  = "RSCLI_PLUGINS"
	pluginsFlag = "plugins"

	openUI      = "ui"
	debugPlugin = "debug"
)

type Plugin interface {
	GetName() string
	Run(args []string) error
}

type PluginWithUi interface {
	GetName() string
	Run(elem uikit.UIElement) uikit.UIElement
}

func Run(args []string) {
	if len(args) == 0 {
		println(help.Run())
		return
	}
	flags := flag.ParseArgs(args)

	if _, ok := flags[debugPlugin]; ok {
		RunDebug(args)
		return
	}

	err := fetchPlugins(flags)
	if err != nil {
		println(help.Header + "error fetching plugins: " + err.Error())
		return
	}

	switch {
	case flags[openUI] != nil:
		RunUI(args)
	default:
		RunCMD(args)
	}

}

func fetchPlugins(args map[string][]string) error {
	pluginsPath, err := findPluginsDir(args)
	if err != nil {
		return errors.Wrap(err, "couldn't find plugins dir")
	}

	if pluginsPath == "" {
		return errors.New("no plugin directory is provided")
	}

	dirs, err := os.ReadDir(pluginsPath)
	if err != nil {
		return errors.Wrap(err, "error reading plugins directory")
	}

	for _, plugPath := range dirs {
		plugName := plugPath.Name()
		if !strings.HasSuffix(plugName, ".so") {
			continue
		}
		var p *plugin.Plugin
		p, err = plugin.Open(path.Join(pluginsPath, plugName))
		if err != nil {
			return errors.Wrapf(err, "error opening plugin %s", plugName)
		}

		var plug plugin.Symbol
		plug, err = p.Lookup("Plug")
		if err != nil {
			return errors.Wrapf(err, "couldn't find \"Plug\" symbol %s", plugName)
		}

		ui, ok := plug.(PluginWithUi)
		if ok {
			pluginsWithUI[ui.GetName()] = ui
			continue
		}

		cli, ok := plug.(Plugin)
		if ok {
			plugins[cli.GetName()] = cli
			continue
		}
		return errors.New("error parsing symbol \"Run\" to any runner in plugin " + plugName)
	}

	return nil
}

func findPluginsDir(args map[string][]string) (string, error) {
	dir, err := flag.ExtractOneValueFromFlags(args, pluginsFlag)
	if err != nil {
		return "", errors.Wrapf(err, "error extracting \"%s\" flag from arguments", pluginsFlag)
	}

	if dir != "" {
		return dir, nil
	}

	pluginsDir, ok := os.LookupEnv(pluginsDIR)
	if !ok {
		return "", nil
	}
	return pluginsDir, nil
}
