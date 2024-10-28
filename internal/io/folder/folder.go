package folder

import (
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/rscli/internal/io"
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
		dir := path.Dir(fl.Name)
		name := path.Base(fl.Name)

		if dir == "." {
			found := false
			for idx := range f.Inner {
				if f.Inner[idx].Name == fl.Name {
					f.Inner[idx] = fl
					found = true
					break
				}
			}
			if !found {
				f.Inner = append(f.Inner, fl)
			}
		} else {
			fl.Name = name
			f.addWithPath(dir, fl)
		}
	}
}

func (f *Folder) addWithPath(pth string, folders ...*Folder) {
	pths := strings.Split(pth, string(os.PathSeparator))
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
				if len(currentFolder.Inner[idx].Content) != 0 && len(folderToAdd.Content) != 0 {
					currentFolder.Inner[idx].Content = folderToAdd.Content
				} else {
					currentFolder.Inner[idx] = folderToAdd
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
	splitPath := make([]string, 0, len(pth))
	for _, p := range pth {
		sp := strings.Split(p, string(os.PathSeparator))
		splitPath = append(splitPath, sp...)
	}

	for _, p := range splitPath {
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

func (f *Folder) Build(rootFolder string) error {
	return f.build(rootFolder)
}

func (f *Folder) build(root string) error {
	pth := path.Join(root, path.Base(f.Name))

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
			err := io.OverrideFile(pth, f.Content)
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
	if f == nil {
		return
	}

	f.isToBeDeleted = true
}

func (f *Folder) Copy() *Folder {
	newF := Folder{
		Name:    f.Name,
		Content: make([]byte, len(f.Content)),
	}

	copy(newF.Content, f.Content)

	return &newF
}

func (f *Folder) CopyWithNewName(name string) *Folder {
	newF := Folder{
		Name:    name,
		Content: make([]byte, len(f.Content)),
	}

	copy(newF.Content, f.Content)

	return &newF
}
