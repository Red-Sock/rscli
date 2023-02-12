package environment

import (
	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli/plugins/environment/scripts"
	"github.com/pkg/errors"
	"os"
)

const PluginName = "environment"

const (
	optionCreate = "create"
	setUp        = "set up environment"
)

const (
	all        = "all"
	selectProj = "select"
)

type handler struct {
	prev uikit.UIElement
}

func Run(prevScreen uikit.UIElement) uikit.UIElement {

	h := handler{prev: prevScreen}

	return radioselect.New(
		h.mainMenuCallback,
		radioselect.Header("Environment"),
		radioselect.Items(
			optionCreate,
			setUp,
		),
	)
}

func (h *handler) mainMenuCallback(arg string) uikit.UIElement {
	switch arg {
	case optionCreate:
		err := scripts.RunCreate()
		if err != nil {
			return h.handleError(err)
		}

		return h.prev
	case setUp:
		return radioselect.New(
			h.handleSetUpSelectEnvs,
			radioselect.Items(
				all,
				selectProj,
			),
		)

	}

	return nil
}

func (h *handler) handleError(err error) uikit.UIElement {
	return label.New(err.Error(), label.NextScreen(func() uikit.UIElement {
		return h.prev
	}))
}

func (h *handler) handleSetUpSelectEnvs(item string) uikit.UIElement {
	dirs, err := os.ReadDir("./")
	if err != nil {
		return h.handleError(err)
	}

	menuItems := make([]string, 0, len(dirs))
	for _, d := range dirs {
		if d.IsDir() && d.Name() != scripts.EnvDir {
			menuItems = append(menuItems, d.Name())
		}
	}

	if len(menuItems) == 0 {
		return h.handleError(errors.New("no project imported to set up"))
	}

	menuItems = append([]string{}, menuItems...)

	switch item {
	case all:
		err = scripts.RunSetUp(menuItems)
		if err != nil {
			return h.handleError(err)
		}
		return h.prev

	case selectProj:
		return multiselect.New(
			h.selectEnvironmentsToSetUp,
			multiselect.Items(menuItems...),
		)
	}
	return h.prev
}

func (h *handler) selectEnvironmentsToSetUp(envs []string) uikit.UIElement {
	err := scripts.RunSetUp(envs)
	if err != nil {
		return h.handleError(err)
	}
	return h.prev
}
