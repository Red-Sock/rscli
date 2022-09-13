package config

import (
	"fmt"
	"github.com/Red-Sock/rscli/internal/commands"
	"github.com/Red-Sock/rscli/internal/service/help"
	"github.com/Red-Sock/rscli/internal/utils"
	"strings"
)

type configTool struct {
}

func NewConfigTool() *configTool {
	return &configTool{}
}

func (c *configTool) Run(args []string) string {
	if utils.Contains(args, commands.Help) {
		return c.HelpMessage()
	}

	args = args[1:]

	if len(args) == 0 {
		return c.runDefault()
	}

	var opts map[string][]string
	var err error

	opts, err = c.parseArgs(args)
	if err != nil {
		return err.Error()
	}

	return NewConfig(opts)
}

func (c *configTool) HelpMessage() string {
	return help.FormMessage(defaultHelp)
}

func (c *configTool) runDefault() string {
	opts := map[string][]string{
		dbFlag: {"postgres", "1"},
	}

	return NewConfig(opts)
}

func (c *configTool) parseArgs(args []string) (map[string][]string, error) {
	flagToArgs := make(map[string][]string)

	key := ""

	for _, item := range args {
		if strings.HasPrefix(item, "-") {
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
--db [source_name]_[connection_name]: allows to setup multiple connections to database. 
        Example: "rscli config --db postgres_users redis_cache" will create configuration file
                 with postgres connection named "users" and redis connection named "cache".
                 
                 "rscli config --db postgres redis" will create configuration file
                 with connections named by source type 
`
