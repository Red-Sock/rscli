package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io/folder"
)

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		if idx := bytes.Index(f.Content, []byte(envpatterns.ProjNamePattern)); idx != -1 {
			f.Content = bytes.ReplaceAll(f.Content, []byte(envpatterns.ProjNamePattern), []byte(name))
			return
		}
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
		if f.Name == envpatterns.ProjNamePattern && len(f.Inner) == 0 {
			f.Name = name
		}
	}
}
