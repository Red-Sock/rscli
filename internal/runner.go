package internal

import (
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

	var err error
	switch {
	case flgs[openUI] != nil:
		delete(flgs, openUI)
		err = RunUI(flgs)
	default:
		err = RunCMD(flgs)
	}

	return err
}
