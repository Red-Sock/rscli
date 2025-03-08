package folder

import (
	"os"
	"path"

	"github.com/Red-Sock/rscli/internal/utils/slices"
)

var ignoredFolders = []string{
	".idea",
	".git",
	".DS_Store",
	"node_modules",
}

func Load(root string) (*Folder, error) {
	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	f := &Folder{
		Name:  root,
		Inner: make([]*Folder, 0, len(dir)),
	}
	for _, d := range dir {
		innerFolder, err := load(path.Join(root, d.Name()), "")
		if err != nil {
			return nil, err
		}
		if innerFolder != nil {
			f.Inner = append(f.Inner, innerFolder)
		}
	}

	return f, nil
}

func load(root, parent string) (*Folder, error) {
	if slices.Contains(ignoredFolders, path.Base(root)) {
		return nil, nil
	}

	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		var innerFile []byte
		innerFile, err = os.ReadFile(root)
		folder := &Folder{
			Name:    path.Base(root),
			Content: innerFile,
		}

		folder.olderVersion = make([]byte, len(folder.Content))
		copy(folder.olderVersion, folder.Content)

		return folder, nil
	}

	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	f := &Folder{
		Name:  path.Base(root),
		Inner: make([]*Folder, 0, len(dir)),
	}

	for _, d := range dir {
		var innerDir *Folder
		innerDir, err = load(path.Join(root, d.Name()), path.Join(parent, f.Name))
		if err != nil {
			return nil, err
		}

		if innerDir != nil {
			f.Inner = append(f.Inner, innerDir)
		}
	}

	return f, nil
}
