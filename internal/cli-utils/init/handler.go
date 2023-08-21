package init

import (
	"github.com/Red-Sock/rscli/internal/stdio"
)

type Handler struct {
	io stdio.IO
}

func (u *Handler) Do(args []string) {
	if len(args) == 0 {
		u.io.Println(GetHelpMessage())
		return
	}

	switch args[0] {

	}
}
