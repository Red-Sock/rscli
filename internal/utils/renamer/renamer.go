package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"
)

func ReplaceProjectName(src []byte, newName string) []byte {
	b := make([]byte, len(src))
	copy(b, src)

	bigName := bytes.ReplaceAll(bytes.ToUpper([]byte(newName)), []byte{'-'}, []byte{'_'})
	b = bytes.ReplaceAll(
		b,
		[]byte(envpatterns.ProjNameCapsPattern),
		bigName,
	)
	smallName := bytes.ReplaceAll(bytes.ToLower([]byte(newName)), []byte{'-'}, []byte{'_'})
	b = bytes.ReplaceAll(b,
		[]byte(envpatterns.ProjNamePattern),
		smallName,
	)

	return b
}

func ReplaceProjectNameStr(src string, newName string) string {
	return string(ReplaceProjectName([]byte(src), newName))
}
