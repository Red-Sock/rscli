package utils

import (
	"fmt"
	"strings"
)

func ParseArgs(args []string) (map[string][]string, error) {
	flagToArgs := make(map[string][]string)

	key := ""

	for _, item := range args {
		if strings.HasPrefix(item, "-") {
			if _, ok := flagToArgs[item]; ok {
				return nil, fmt.Errorf("%s flag repited", item)
			}
			key = strings.ReplaceAll(item, "-", "")
			flagToArgs[key] = nil
		} else {
			flagToArgs[key] = append(flagToArgs[key], item)
		}
	}

	if emptyArgs, ok := flagToArgs[""]; ok {
		return nil, fmt.Errorf("unknown arguments %s", strings.Join(emptyArgs, ";"))
	}

	return flagToArgs, nil
}
