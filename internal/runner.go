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

const (
	pluginExtension = ".so"
)

type Plugin interface {
	GetName() string
	Run(args map[string][]string) error
}

type PluginWithUi interface {
	GetName() string
	Run(elem uikit.UIElement) uikit.UIElement
}

func Run(args []string) error {

	if len(args) == 0 {
		println(help.Run())
		return nil
	}

	flgs := flag.ParseArgs(args)

	ok, err := execIfEmbededCommand(flgs)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	err = fetchPlugins(flgs)
	if err != nil {
		return errors.New(help.Header + "error fetching plugins: " + err.Error())
	}

	switch {
	case flgs[openUI] != nil:
		delete(flgs, openUI)
		err = RunUI(flgs)
	default:
		err = RunCMD(flgs)
	}

	return err
}

func execIfEmbededCommand(flags map[string][]string) (ok bool, err error) {
	for _, b := range basicPlugin {
		if _, ok = flags[b.GetName()]; ok {

			delete(flags, b.GetName())
			err = b.Run(flags)
			if err != nil {
				return ok, err
			}
		}
	}

	return ok, nil
}

func fetchPlugins(args map[string][]string) error {
	pluginsPath, err := findPluginsDir(args)
	if err != nil {
		return errors.Wrap(err, "couldn't find plugins dir")
	}

	if pluginsPath == "" {
		return fmt.Errorf("plugin directory doesn't exist. %s %s to fix", commands.RsCLI(), commands.FixUtil)
	}

	err = scanPluginsDir(pluginsPath)
	if err != nil {
		return errors.Wrap(err, "error scanning plugins dir")
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

func scanPluginsDir(pluginsPath string) error {
	dirs, err := os.ReadDir(pluginsPath)
	if err != nil {
		return errors.Wrap(err, "error reading plugins directory")
	}

	for _, item := range dirs {
		tmpPath := path.Join(pluginsPath, item.Name())

		if strings.HasSuffix(tmpPath, pluginExtension) {
			err = fetchPlugin(tmpPath)
		} else {
			err = scanPluginsDir(tmpPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func fetchPlugin(plugPath string) (err error) {
	if !strings.HasSuffix(plugPath, ".so") {
		return nil
	}

	var p *plugin.Plugin
	p, err = plugin.Open(plugPath)
	if err != nil {
		return errors.Wrapf(err, "error opening plugin")
	}

	var plug plugin.Symbol
	plug, err = p.Lookup("Plug")
	if err != nil {
		return errors.Wrapf(err, "couldn't find \"Plug\" symbol %s", plugPath)
	}

	ui, ok := plug.(PluginWithUi)
	if ok {
		pluginsWithUI[ui.GetName()] = ui
		return nil
	}

	cli, ok := plug.(Plugin)
	if ok {
		plugins[cli.GetName()] = cli
		return nil
	}

	return errors.New("error parsing symbol \"Run\" to any runner in plugin " + plugPath)
}
