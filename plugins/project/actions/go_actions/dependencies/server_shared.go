package dependencies

import (
	"path"

	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

func initServerFiles(proj Project) {
	serverManagerPath := []string{patterns.InternalFolder, patterns.TransportFolder, patterns.ServerManager.Name}

	if proj.GetFolder().GetByPath(serverManagerPath...) == nil {
		proj.GetFolder().Add(
			patterns.ServerManager.
				CopyWithNewName(path.Join(serverManagerPath...)))
	}
}
