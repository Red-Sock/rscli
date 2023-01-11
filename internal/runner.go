package internal

import (
	"fmt"
	"github.com/Red-Sock/rscli/pkg/commands"
	"github.com/Red-Sock/rscli/pkg/flag/flags"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/pkg/errors"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/pkg/service/help"
)

const (
	openUI = "ui"
)

type Plugin interface {
	GetName() string
	Run(args map[string][]string) error
}

type PluginWithUi interface {
	GetName() string
	Run(elem uikit.UIElement) uikit.UIElement
}

func Run(args []string) {
	if len(args) == 1 {
		println(help.Run())
		os.Exit(0)
	}
	flgs := flag.ParseArgs(args)

	err := ifBasicCommands(flgs)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}

	err = fetchPlugins(flgs)
	if err != nil {
		println(help.Header + "error fetching plugins: " + err.Error())
		return
	}

	switch {
	case flgs[openUI] != nil:
		RunUI(flgs)
	default:
		RunCMD(flgs)
	}

}

func ifBasicCommands(flags map[string][]string) error {
	for _, b := range basicPlugin {
		if _, ok := flags[b.GetName()]; ok {

			delete(flags, b.GetName())
			err := b.Run(flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fetchPlugins(args map[string][]string) error {
	pluginsPath, err := findPluginsDir(args)
	if err != nil {
		return errors.Wrap(err, "couldn't find plugins dir")
	}

	if pluginsPath == "" {
		return fmt.Errorf("plugin directory doesn't exist. %s %s fix to fix", commands.RsCLI, commands.FixUtil)
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
	dir, err := flag.ExtractOneValueFromFlags(args, flags.PluginsDirFlag)
	if err != nil {
		return "", errors.Wrapf(err, "error extracting \"%s\" flag from arguments", flags.PluginsDirFlag)
	}

	if dir != "" {
		return dir, nil
	}

	pluginsDir, ok := os.LookupEnv(flags.PluginsDirEnv)
	if !ok {
		return "", nil
	}
	return pluginsDir, nil
}
