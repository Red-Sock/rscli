package config

import (
	"fmt"
	"strings"
)

func NewConfig(opts map[string][]string) string {
	parts := make([]fmt.Stringer, 0, len(opts))

	for key, args := range opts {

		flag := strings.Replace(key, "-", "", -1)
		p, ok := patterns[flag]
		if !ok {
			return fmt.Sprintf("Unknown flag %s. Use \"rscli config help\" for help", key)
		}
		parts = append(parts, p(args).getParts(0)...)

	}
	return ""
}
