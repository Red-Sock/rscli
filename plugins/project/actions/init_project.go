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
		go_actions.PrepareProjectStructureAction{}, // basic go project structure
		go_actions.PrepareGoConfigFolderAction{},   // generates config keys

		go_actions.BuildProjectAction{}, // build project in file system

		go_actions.InitGoModAction{},       // executes go mod
		go_actions.PrepareMakefileAction{}, // prepare Makefile

		go_actions.BuildProjectAction{}, // builds project to file system

		go_actions.RunGoTidyAction{}, // adds/clears project initialization(api, resources) and replaces project name template with actual project name
		go_actions.RunGoFmtAction{},  // fetches dependencies and formats go code

		git.InitGit{},
	}
}
