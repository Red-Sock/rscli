package folder

import (
	"os"
	"path"
)

type Folder struct {
	Name    string
	Inner   []*Folder
	Content []byte
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

	folder.Inner = append(folder.Inner, folders...)
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
		fw, err := os.Create(pth)
		if err != nil {
			return err
		}
		defer fw.Close()
		_, err = fw.Write(f.Content)
		return err
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
