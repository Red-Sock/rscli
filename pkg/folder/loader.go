package folder

import (
	"os"
	"path"
)

func Load(pth string) (Folder, error) {
	f := Folder{}

	_, f.Name = path.Split(pth)

	dir, err := os.ReadDir(pth)
	if err != nil {
		return f, err
	}

	for _, d := range dir {
		if d.IsDir() {
			var innerDir Folder
			innerDir, err = Load(path.Join(pth, d.Name()))
			if err != nil {
				return Folder{}, err
			}

			f.Inner = append(f.Inner, &innerDir)
		} else {
			fileName := d.Name()

			var innerFile []byte
			innerFile, err = os.ReadFile(path.Join(pth, fileName))
			folder := &Folder{
				Name:    fileName,
				Content: innerFile,
			}
			folder.olderVersion = make([]byte, len(folder.Content))
			copy(folder.olderVersion, folder.Content)
			f.Inner = append(f.Inner, folder)
		}
	}

	return f, nil
}
