package internal

import (
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/service/help"
	"github.com/Red-Sock/rscli/internal/utils"
)

type Tool interface {
	Run(args []string) (output string)
	HelpMessage() string
}

func Run(args []string) {
	if len(args) == 0 {
		println(help.Run())
		return
	}
	var res string
	switch {
	case utils.Contains(config.Command(), args[0]):
		res = config.Run(args[0:])
	}

	println(res)
}
