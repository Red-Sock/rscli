package flag

import "strings"

func MergeFlags(flags map[string][]string) string {
	sb := &strings.Builder{}
	for f, v := range flags {
		sb.WriteString(f + " " + strings.Join(v, " "))
	}
	return sb.String()
}
