package grpc_discovery

import (
	"os"
	"path"
	"sort"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
)

var modFolderPath = os.Getenv("GOPATH") + "/pkg/mod/"

func GetPathToGlobalModule(packageName string) (pathToModule string, err error) {
	packageName = FilterPackageName(packageName)
	packagePath := path.Join(modFolderPath, FilterPackageName(packageName))

	if strings.Contains(path.Base(packagePath), "@") {
		return packagePath, nil
	}

	root := path.Dir(packagePath)
	potentialDirs, err := os.ReadDir(root)
	if err != nil {
		return "", errors.Wrap(err, "error reading potential packages paths")
	}

	baseName := path.Base(packageName)
	moveIdx := 0
	for idx := range potentialDirs {
		if !strings.HasPrefix(potentialDirs[idx].Name(), baseName) {
			potentialDirs[moveIdx], potentialDirs[idx] = potentialDirs[idx], potentialDirs[moveIdx]
			moveIdx++
		}

	}
	potentialDirs = potentialDirs[moveIdx:]

	sort.Slice(potentialDirs, func(i, j int) bool {
		return potentialDirs[i].Name() < potentialDirs[i].Name()
	})

	packagePath = path.Join(root, potentialDirs[0].Name())

	return packagePath, nil
}
