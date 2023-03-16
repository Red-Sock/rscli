package project

import (
	"os"

	uikit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/basic/label"

	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	"github.com/Red-Sock/rscli/plugins/project/processor"
)

func Tidy(prev uikit.UIElement) uikit.UIElement {
	wd, err := os.Getwd()
	if err != nil {
		return shared_ui.GetHeaderFromText(err.Error())
	}

	err = processor.Tidy(wd)
	if err != nil {
		return shared_ui.GetHeaderFromLabel(label.New(err.Error(), label.NextScreen(prev)))
	}

	return prev
}
