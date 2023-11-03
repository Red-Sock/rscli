package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

var (
	ErrNoFolderInConfig = errors.New("no folder path in rscli config")
)

// containsDependencyFolder - searches through RSCLI_PATH_TO_CLIENTS
// folders in order to find depName
// IF dependency already placed - returns path to it
func containsDependencyFolder(paths []string, rootF *folder.Folder, depName string) (ok bool, err error) {
	if len(paths) == 0 {
		return false, errors.Wrap(ErrNoFolderInConfig, "no client")
	}

	for _, clientPath := range paths {
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
