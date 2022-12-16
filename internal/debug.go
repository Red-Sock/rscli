package internal

import (
	rscliuitkit "github.com/Red-Sock/rscli-uikit"
	ui "github.com/Red-Sock/rscli/plugins/src/config/ui/manager"
)

// Here you suppose to call a debugging method

func RunDebug(args []string) {
	args = args[1:]

	q := make(chan struct{})
	rscliuitkit.NewHandler(ui.Run(nil)).Start(q)
}
