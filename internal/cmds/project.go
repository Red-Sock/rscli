package cmds

import (
	"github.com/Red-Sock/rscli/internal/service/project"
)

func RunProject(args []string) {
	p, err := project.NewProject(args)
	if err != nil {
		println(err.Error())
		return
	}

	println(p)
}
