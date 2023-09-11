package renamer

import (
	"bytes"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
)

func ReplaceProjectName(src []byte, newName string) []byte {
	b := make([]byte, len(src))
	copy(b, src)

	bigName := bytes.ReplaceAll(bytes.ToUpper([]byte(newName)), []byte{'-'}, []byte{'_'})
	b = bytes.ReplaceAll(
		b,
		[]byte(patterns.ProjNameCapsPattern),
		bigName,
	)
	smallName := bytes.ReplaceAll(bytes.ToLower([]byte(newName)), []byte{'-'}, []byte{'_'})
	b = bytes.ReplaceAll(b,
		[]byte(patterns.ProjNamePattern),
		smallName,
	)
	return b
}

func ReplaceProjectNameStr(src string, newName string) []byte {
	return ReplaceProjectName([]byte(src), newName)
}
