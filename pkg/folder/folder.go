package folder

import (
	"os"
	"path"
	"strings"
)

type Folder struct {
	Name    string
	Inner   []*Folder
	Content []byte

	olderVersion  []byte
	isToBeDeleted bool
}

func (f *Folder) Add(folders ...*Folder) {
	for _, fl := range folders {
		if splited := strings.Split(fl.Name, string(os.PathSeparator)); len(splited) > 1 {
			fl.Name = splited[len(splited)-1]
			f.AddWithPath(splited[:len(splited)-1], fl)
		} else {
			f.Inner = append(f.Inner, fl)
		}
	}
}

func (f *Folder) AddWithPath(pths []string, folders ...*Folder) {
	f.addWithPath(pths, false, folders)
}

func (f *Folder) ForceAddWithPath(pths []string, folders ...*Folder) {
	f.addWithPath(pths, true, folders)
}

func (f *Folder) addWithPath(pths []string, needToReplace bool, folders []*Folder) {
	if len(folders) == 0 {
		return
	}

	currentFolder := f
	for _, pathPart := range pths {
		var pathFolder *Folder

		for currentFolderIdx := range currentFolder.Inner {
			if currentFolder.Inner[currentFolderIdx].Name == pathPart {
				pathFolder = currentFolder.Inner[currentFolderIdx]
				break
			}
		}
		if pathFolder == nil {
			pathFolder = &Folder{
				Name: pathPart,
			}
			currentFolder.Inner = append(currentFolder.Inner, pathFolder)
		}
		currentFolder = pathFolder
	}
	for _, folderToAdd := range folders {
		var isAdded bool

		for idx, itemInCurrentFolder := range currentFolder.Inner {
			if itemInCurrentFolder.Name == folderToAdd.Name {
				if needToReplace {
					if len(currentFolder.Inner[idx].Content) != 0 && len(folderToAdd.Content) != 0 {
						currentFolder.Inner[idx].Content = folderToAdd.Content
					} else {
						currentFolder.Inner[idx] = folderToAdd
					}
				}
				isAdded = true
				break
			}
		}

		if !isAdded {
			currentFolder.Inner = append(currentFolder.Inner, folderToAdd)
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

func (f *Folder) Build() error {
	return f.build("")
}

func (f *Folder) build(root string) error {

	pth := path.Join(root, f.Name)

	if f.isToBeDeleted {
		err := os.RemoveAll(pth)
		if err != nil {
			return err
		}
		return nil
	}

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

		if len(f.Content) != 0 && !(len(f.Content) == 1 && f.Content[0] != 0) {
			err := os.WriteFile(pth, f.Content, 0755)
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
		err = d.build(pth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Folder) Delete() {
	f.isToBeDeleted = true
}
