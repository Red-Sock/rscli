package renamer

import (
	"bytes"
	"strings"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
)

func ReplaceProjectName(src []byte, newName string) []byte {
	b := make([]byte, len(src))
	copy(b, src)
	b = bytes.ReplaceAll(b, []byte(patterns.ProjNameCapsPattern), []byte(strings.ToUpper(newName)))
	b = bytes.ReplaceAll(b, []byte(patterns.ProjNamePattern), []byte(strings.ToLower(newName)))
	return b
}
