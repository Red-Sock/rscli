package config

import (
	"errors"
	"fmt"
	"github.com/Red-Sock/rscli/internal/service/help"
	"github.com/Red-Sock/rscli/internal/utils"
	"strings"
)

var commands = []string{"config", "c"}

func Command() []string {
	return commands
}

func Run(args []string) (*Config, error) {
	if utils.Contains(args, help.Command) {
		return nil, HelpMessage()
	}

	args = args[1:]

	if len(args) == 0 {
		return runDefault()
	}

	var opts map[string][]string
	var err error

	opts, err = parseArgs(args)
	if err != nil {
		return nil, err
	}

	return NewConfig(opts)
}

func HelpMessage() error {
	return errors.New(help.FormMessage(defaultHelp))
}

func runDefault() (*Config, error) {
	opts := map[string][]string{
		sourceNamePg: {"postgres"},
	}

	return NewConfig(opts)
}

func parseArgs(args []string) (map[string][]string, error) {
	flagToArgs := make(map[string][]string)

	key := ""

	for _, item := range args {
		if strings.HasPrefix(item, "-") {
			if _, ok := flagToArgs[item]; ok {
				return nil, fmt.Errorf("%s flag repited", item)
			}
			key = item
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

const defaultHelp = `
rscli config - creates configuration file
============================================================================
--pg [connection_name]: setup postgres connection(s).
        Example: "rscli config --pg" will create configuration file
                 with default connection to local postgres.

--rds [connection_name]: setup redis connection(s).
       Example:  "rscli config --rds" will create configuration file
                 with default connection to local redis.

Putting name(s) after flag will create named connections
       Example:  "rscli config --pg parking_lot key" will create
                 configuration file with connections to postgres
                 databases named "parking_lot" and "key". 

--fo forces to override currently available config.

--path [path/to/config.yml] will set custom config path. 
`
