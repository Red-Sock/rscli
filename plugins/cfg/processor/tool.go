package processor

import (
	"errors"
)

var commands = []string{"config", "c"}

func Command() []string {
	return commands
}

func Run(args map[string][]string) (*Config, error) {
	if _, ok := args["help"]; ok {
		return nil, errors.New(defaultHelp)
	}

	if len(args) == 0 {
		return runDefault()
	}

	return NewConfig(args)
}

func runDefault() (*Config, error) {
	opts := map[string][]string{}

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
