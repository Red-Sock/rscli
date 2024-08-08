package actions

import (
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/pipelines"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
)

func GetTidyActionsForProject(pt proj_interfaces.ProjectType) []Action {
	switch pt {
	case proj_interfaces.ProjectTypeGo:
		return append(
			append(
				commonProjectTidyPreActions(),
				goProjectTidyActions()...),
			commonProjectTidyPostActions()...)
	default:
		return unknownProjectActions()
	}
}

func goProjectTidyActions() []Action {
	return []Action{
		go_actions.PrepareGoConfigFolderAction{},
		go_actions.PrepareMakefileAction{},
		go_actions.PrepareClientsAction{},
		go_actions.BuildProjectAction{},
		go_actions.RunMakeGenAction{},
		go_actions.BuildProjectAction{},
		go_actions.RunGoTidyAction{},
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
