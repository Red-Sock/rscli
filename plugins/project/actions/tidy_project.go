package actions

import (
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/pipelines"
)

func GetTidyActionsForProject(pt project.Type) []Action {
	out := commonProjectTidyPreActions()

	switch pt {
	case project.TypeGo:
		out = append(out, goProjectTidyActions()...)
	default:
		return unknownProjectActions()
	}

	out = append(out, commonProjectTidyPostActions()...)
	return out
}

func goProjectTidyActions() []Action {
	return []Action{
		go_actions.PrepareConfigFolder{},
		go_actions.PrepareMakefile{},
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},
		go_actions.PrepareDockerfile{},
		go_actions.BuildProjectAction{},
		// Takes too long
		//go_actions.RunMakeGenAction{},
		go_actions.InitGoProjectApp{},

		go_actions.BuildProjectAction{},
		// Takes too long
		// go_actions.RunGoTidyAction{},
		go_actions.GoFmt{},
	}
}

func commonProjectTidyPreActions() []Action {
	return []Action{
		pipelines.TidyGithubWorkflowAction{},
	}
}

func commonProjectTidyPostActions() []Action {
	return []Action{
		git.CommitWithUntrackedAction{},
	}
}

func unknownProjectActions() []Action {
	return nil
}
