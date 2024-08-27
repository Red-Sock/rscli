package actions

import (
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/dockerfile_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/pipelines"
)

func GetTidyActionsForProject(pt project.ProjectType) []Action {
	out := commonProjectTidyPreActions()

	switch pt {
	case project.ProjectTypeGo:
		out = append(out, goProjectTidyActions()...)
	default:
		return unknownProjectActions()
	}

	out = append(out, commonProjectTidyPostActions()...)
	return out
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
		go_actions.RunGoFmtAction{},
	}
}

func commonProjectTidyPreActions() []Action {
	return []Action{
		pipelines.TidyGithubWorkflowAction{},
	}
}

func commonProjectTidyPostActions() []Action {
	return []Action{
		dockerfile_actions.DockerFileTidyAction{},
		git.CommitWithUntrackedAction{},
	}
}

func unknownProjectActions() []Action {
	return nil
}
