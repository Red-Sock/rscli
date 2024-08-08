package pipelines

import (
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type TidyGithubWorkflowAction struct {
}

func (a TidyGithubWorkflowAction) Do(p proj_interfaces.Project) error {
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
	case proj_interfaces.ProjectTypeGo:
		ghF.Add(projpatterns.GithubWorkflowGoBranchPush.Copy())
	}

	return nil
}
func (a TidyGithubWorkflowAction) NameInAction() string {
	return "Tiding github workflows"
}
