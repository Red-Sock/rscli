package project

import (
	"os"
	"path"
	"strings"
)

type Folder struct {
	name    string
	inner   []*Folder
	content []byte
}

func (f *Folder) AddWithPath(pth string, folders ...*Folder) {
	if len(folders) == 0 {
		return
	}
	pths := strings.Split(pth, string(os.PathSeparator))

	folder := f
	for _, pathPart := range pths {
		var pathFolder *Folder

		for currentFolderIdx := range folder.inner {
			if folder.inner[currentFolderIdx].name == pathPart {
				pathFolder = folder.inner[currentFolderIdx]
				break
			}
		}
		if pathFolder == nil {
			pathFolder = &Folder{
				name: pathPart,
			}
			folder.inner = append(folder.inner, pathFolder)
		}
		folder = pathFolder
	}

	folder.inner = append(folder.inner, folders...)
}

func (f *Folder) Build(root string) error {
	pth := path.Join(root, f.name)

	if len(f.content) != 0 {
		fw, err := os.Create(pth)
		if err != nil {
			return err
		}
		defer fw.Close()
		_, err = fw.Write(f.content)
		return err
	}

	err := os.MkdirAll(pth, 0755)
	if err != nil {
		return err
	}

	for _, d := range f.inner {
		err = d.Build(pth)
		if err != nil {
			return err
		}
	}

	return nil
}
