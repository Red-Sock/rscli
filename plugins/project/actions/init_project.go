package actions

import (
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
)

func InitProject(pt project.Type) []Action {
	switch pt {
	case project.TypeGo:
		return initVirtualGoProject()
	default:
		return nil
	}
}

func initVirtualGoProject() []Action {
	return []Action{
		go_actions.PrepareProjectStructure{}, // basic go project structure
		go_actions.InitGoProjectApp{},
		go_actions.PrepareConfigFolder{}, // generates config keys
		go_actions.PrepareMakefile{},
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},

		go_actions.BuildProjectAction{}, // build project in file system

		go_actions.InitGoMod{}, // executes go mod

		go_actions.BuildProjectAction{}, // builds project to file system

		go_actions.RunGoTidyAction{}, // adds/clears project initialization(api, resources) and replaces project name template with actual project name
		go_actions.GoFmt{},           // fetches dependencies and formats go code

		git.InitGit{},
	}
}
