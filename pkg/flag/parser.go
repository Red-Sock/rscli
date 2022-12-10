package flag

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNoArgumentsSpecifiedForFlag = errors.New("flag specified but no name was given")
	ErrFlagHasTooManyArguments     = errors.New("too many arguments specified for flag")
)

func ExtractOneValueFromFlags(flagsArgs map[string][]string, flags ...string) (string, error) {
	var name []string
	for _, f := range flags {
		var ok bool
		name, ok = flagsArgs[f]
		if ok {
			break
		}
	}

	if name == nil {
		return "", nil
	}

	if len(name) == 0 {
		return "", fmt.Errorf("%w expected 1 got 0 ", ErrNoArgumentsSpecifiedForFlag)
	}

	if len(name) > 1 {
		return "", fmt.Errorf("%w expected 1 got %d", ErrFlagHasTooManyArguments, len(name))
	}

	return name[0], nil
}

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
