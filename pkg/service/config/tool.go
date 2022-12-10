package config

import (
	"errors"

	"github.com/Red-Sock/rscli/pkg/flag"

	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/pkg/service/help"
)

var commands = []string{"config", "c"}

func Command() []string {
	return commands
}

func Run(args []string) (*Config, error) {
	if slices.Contains(args, help.Command) {
		return nil, HelpMessage()
	}

	if len(args) == 0 {
		return runDefault()
	}

	var opts map[string][]string
	var err error

	opts, err = flag.ParseArgs(args)
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
		SourceNamePg: {"postgres"},
	}

	return NewConfig(opts)
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
