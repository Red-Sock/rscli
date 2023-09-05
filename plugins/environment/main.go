package environment

import (
	"os"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	"github.com/Red-Sock/rscli-uikit/composit-items/radioselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"
	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/patterns"
	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	"github.com/Red-Sock/rscli/plugins/environment/scripts"
)

const PluginName = "environment"

const (
	setUp = "set up environment"
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
		radioselect.HeaderLabel(shared_ui.GetHeaderFromText("Environment")),
		radioselect.Items(
			setUp,
		),
		radioselect.Position(common.NewRelativePositioning(common.NewFillSpacePositioning(), common.NewFillSpacePositioning(), 0.4, 0.4)),
	)
}

func (h *handler) mainMenuCallback(arg string) uikit.UIElement {
	switch arg {

	case setUp:
		return radioselect.New(
			h.handleSetUpSelectEnvs,
			radioselect.HeaderLabel(shared_ui.GetHeaderFromText("What environment to set up ?")),
			radioselect.Items(
				all,
				selectProj,
			),
			radioselect.Position(common.NewRelativePositioning(common.NewFillSpacePositioning(), common.NewFillSpacePositioning(), 0.4, 0.4)),
		)
	}

	return nil
}

func (h *handler) handleError(err error) uikit.UIElement {
	return shared_ui.GetHeaderFromLabel(
		label.New(err.Error(),
			label.NextScreen(h.prev)))
}

func (h *handler) handleSetUpSelectEnvs(item string) uikit.UIElement {
	dirs, err := os.ReadDir("./")
	if err != nil {
		return h.handleError(err)
	}

	menuItems := make([]string, 0, len(dirs))
	for _, d := range dirs {
		if d.IsDir() && d.Name() != patterns.EnvDir {
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
			multiselect.HeaderLabel(shared_ui.GetHeaderFromText("Found this RSCLI compatible projects:")),
			multiselect.Items(menuItems...),
			multiselect.Position(common.NewRelativePositioning(common.NewFillSpacePositioning(), common.NewFillSpacePositioning(), 0.4, 0.4)),
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
