package folder

import (
	"os"
	"path"
)

type Folder struct {
	Name    string
	Inner   []*Folder
	Content []byte

	olderVersion []byte
}

func (f *Folder) Add(folder ...*Folder) {
	f.Inner = append(f.Inner, folder...)
}

func (f *Folder) AddWithPath(pths []string, folders ...*Folder) {
	if len(folders) == 0 {
		return
	}

	folder := f
	for _, pathPart := range pths {
		var pathFolder *Folder

		for currentFolderIdx := range folder.Inner {
			if folder.Inner[currentFolderIdx].Name == pathPart {
				pathFolder = folder.Inner[currentFolderIdx]
				break
			}
		}
		if pathFolder == nil {
			pathFolder = &Folder{
				Name: pathPart,
			}
			folder.Inner = append(folder.Inner, pathFolder)
		}
		folder = pathFolder
	}
	for _, folderToAdd := range folders {
		var isAdded bool

		for idx, existingItem := range folder.Inner {
			if existingItem.Name == folderToAdd.Name {
				folder.Inner[idx] = folderToAdd
				isAdded = true
				break
			}
		}

		if !isAdded {
			folder.Inner = append(folder.Inner, folderToAdd)
		}
	}
}

func (f *Folder) GetByPath(pth ...string) *Folder {
	currentFolder := f
	for _, p := range pth {
		var foundFolder *Folder
		for _, cf := range currentFolder.Inner {
			if cf.Name == p {
				foundFolder = cf
				break
			}
		}

		if foundFolder == nil {
			return nil
		}
		currentFolder = foundFolder
	}

	return currentFolder
}

func (f *Folder) Build(root string) error {
	pth := path.Join(root, f.Name)

	if len(f.Content) != 0 {

		if len(f.olderVersion) == len(f.Content) {
			var idx int
			for idx = range f.olderVersion {
				if f.olderVersion[idx] != f.Content[idx] {
					break
				}
			}
			if len(f.olderVersion) != idx-1 {
				return nil
			}
		}
		err := os.RemoveAll(pth)
		if err != nil {
			return err
		}

		if len(f.Content) != 0 && !(len(f.Content) == 1 && f.Content[0] != 0) {
			err = os.WriteFile(pth, f.Content, 0755)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := os.MkdirAll(pth, 0755)
	if err != nil {
		return err
	}

	for _, d := range f.Inner {
		err = d.Build(pth)
		if err != nil {
			return err
		}
	}

	return nil
}
