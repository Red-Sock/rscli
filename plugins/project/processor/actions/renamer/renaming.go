package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		if idx := bytes.Index(f.Content, []byte(patterns.ImportProjectNamePatternKebabCase)); idx != -1 {
			f.Content = bytes.ReplaceAll(f.Content, []byte(patterns.ImportProjectNamePatternKebabCase), []byte(name))
			return
		}
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
		if f.Name == patterns.ImportProjectNamePatternKebabCase && len(f.Inner) == 0 {
			f.Name = name
		}
	}
}
