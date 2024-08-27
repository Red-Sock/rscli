package pipelines

import (
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

type TidyGithubWorkflowAction struct {
}

func (a TidyGithubWorkflowAction) Do(p project.Project) error {
	ghF := p.GetFolder().GetByPath(projpatterns.GithubFolder, projpatterns.WorkflowsFolder)
	if ghF == nil {
		ghF = &folder.Folder{
			Name: projpatterns.GithubFolder,
			Inner: []*folder.Folder{
				{
					Name: projpatterns.WorkflowsFolder,
				},
			},
		}
		p.GetFolder().Add(ghF)
		ghF = ghF.Inner[0]
	}

	ghF.Add(projpatterns.GithubWorkflowRelease.Copy())

	switch p.GetType() {
	case project.ProjectTypeGo:
		ghF.Add(projpatterns.GithubWorkflowGoBranchPush.Copy())
	}

	return nil
}
func (a TidyGithubWorkflowAction) NameInAction() string {
	return "Tiding github workflows"
}
