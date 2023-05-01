package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/pkg/folder"
)

const ImportProjectNamePattern = "financial-microservice"

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		if idx := bytes.Index(f.Content, []byte(ImportProjectNamePattern)); idx != -1 {
			f.Content = bytes.ReplaceAll(f.Content, []byte(ImportProjectNamePattern), []byte(name))
			return
		}
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
		if f.Name == ImportProjectNamePattern && len(f.Inner) == 0 {
			f.Name = name
		}
	}
}
