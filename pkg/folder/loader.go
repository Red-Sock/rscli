package folder

import (
	"os"
	"path"

	"github.com/Red-Sock/rscli/internal/utils/slices"
)

var ignoredFolders = []string{
	".idea",
	".git",
}

func Load(root, projPath string) (Folder, error) {
	f := Folder{}

	pth := path.Join(root, projPath)

	_, f.Name = path.Split(pth)

	dir, err := os.ReadDir(pth)
	if err != nil {
		return f, err
	}

	for _, d := range dir {
		if slices.Contains(ignoredFolders, path.Join(pth, d.Name())) {
			continue
		}

		nodeName := d.Name()

		if slices.Contains(ignoredFolders, path.Join(projPath, nodeName)) {
			continue
		}

		if d.IsDir() {
			var innerDir Folder
			innerDir, err = Load(root, path.Join(projPath, nodeName))
			if err != nil {
				return Folder{}, err
			}

			f.Inner = append(f.Inner, &innerDir)
		} else {

			var innerFile []byte
			innerFile, err = os.ReadFile(path.Join(pth, nodeName))
			folder := &Folder{
				Name:    nodeName,
				Content: innerFile,
			}
			folder.olderVersion = make([]byte, len(folder.Content))
			copy(folder.olderVersion, folder.Content)
			f.Inner = append(f.Inner, folder)
		}
	}

	return f, nil
}
