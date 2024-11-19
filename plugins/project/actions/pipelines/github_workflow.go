package pipelines

import (
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

type TidyGithubWorkflowAction struct {
}

func (a TidyGithubWorkflowAction) Do(p project.IProject) error {
	ghF := p.GetFolder().GetByPath(patterns.GithubFolder, patterns.WorkflowsFolder)
	if ghF == nil {
		ghF = &folder.Folder{
			Name: patterns.GithubFolder,
			Inner: []*folder.Folder{
				{
					Name: patterns.WorkflowsFolder,
				},
			},
		}
		p.GetFolder().Add(ghF)
		ghF = ghF.Inner[0]
	}

	ghF.Add(patterns.GithubWorkflowRelease.Copy())

	switch p.GetType() {
	case project.TypeGo:
		ghF.Add(patterns.GithubWorkflowGoBranchPush.Copy())
	}

	return nil
}
func (a TidyGithubWorkflowAction) NameInAction() string {
	return "Tiding github workflows"
}
