package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/pkg/folder"
)

const ProjectNamePattern = "financial-microservice"

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		if idx := bytes.Index(f.Content, []byte(ProjectNamePattern)); idx != -1 {
			f.Content = bytes.ReplaceAll(f.Content, []byte(ProjectNamePattern), []byte(name))
			return
		}
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
		if f.Name == ProjectNamePattern && len(f.Inner) == 0 {
			f.Name = name
		}
	}
}
