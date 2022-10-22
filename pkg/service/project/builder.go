package project

import (
	"os"
	"path"
)

type folder struct {
	name    string
	inner   []folder
	content []byte
}

func (f *folder) MakeAll(root string) error {
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
		err = d.MakeAll(pth)
		if err != nil {
			return err
		}
	}
	return nil
}
