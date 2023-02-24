package shared

import uikit "github.com/Red-Sock/rscli-uikit"

type Interactor interface {
	Ask(string) (string, error)
	RunScreen(func(prevScreen uikit.UIElement) uikit.UIElement) error
}
