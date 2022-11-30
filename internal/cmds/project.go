package cmds

import (
	"github.com/Red-Sock/rscli/pkg/service/project"
)

func RunProject(args []string) {
	_, err := project.NewProjectWithRowArgs(args)
	if err != nil {
		println(err.Error())
		return
	}
}
