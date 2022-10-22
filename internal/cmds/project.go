package cmds

import (
	"github.com/Red-Sock/rscli/pkg/service/project"
)

func RunProject(args []string) {
	_, err := project.NewProject(args)
	if err != nil {
		println(err.Error())
		return
	}
}
