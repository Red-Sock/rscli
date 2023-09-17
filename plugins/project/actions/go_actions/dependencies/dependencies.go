package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
)

var ErrNoClientFolderInConfig = errors.New("no client path in rscli config")

// containsDependency - searches through RSCLI_PATH_TO_CLIENTS
// folders in order to find depName
// IF dependency already placed - returns path to it
func containsDependency(cfg *rscliconfig.RsCliConfig, rootF *folder.Folder, depName string) (ok bool, err error) {
	if len(cfg.Env.PathsToClients) == 0 {
		return false, errors.Wrap(ErrNoClientFolderInConfig, "no client")
	}

	for _, clientPath := range cfg.Env.PathsToClients {
		clientFolder := rootF.GetByPath(clientPath)
		if clientFolder == nil {
			continue
		}

		for _, cF := range clientFolder.Inner {
			if cF.Name == depName {
				return true, nil
			}
		}
	}

	return false, nil
}
