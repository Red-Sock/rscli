package dependencies

import (
	"path"

	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

func applyServerFolder(proj project.Project) {
	serverManagerPath := []string{projpatterns.InternalFolder, projpatterns.TransportFolder, projpatterns.ServerManager.Name}
	if proj.GetFolder().GetByPath(serverManagerPath...) == nil {
		proj.GetFolder().Add(
			projpatterns.ServerManager.
				CopyWithNewName(path.Join(serverManagerPath...)))
	}
}
