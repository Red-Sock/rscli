package project

import (
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/config"
)

// tries to find path to configuration in same directory
func findConfigPath() (pth string, err error) {
	currentDir := "./"

	var dirs []os.DirEntry
	dirs, err = os.ReadDir(currentDir)
	if err != nil {
		return "", err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = path.Join(currentDir, config.DefaultDir)
			break
		}
	}

	if pth == "" {
		return "", nil
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return "", err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, config.FileName) {
			pth = path.Join(pth, name)
			break
		}
	}

	return pth, nil
}
