package internal

import (
	"fmt"
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/pkg/errors"
)

func RunCMD(args map[string][]string) error {
	if len(args) == 0 {
		return errors.New("no args given")
	}

	for _, plugin := range plugins {
		if _, ok := args[plugin.GetName()]; ok {

			delete(args, plugin.GetName())
			err := plugin.Run(args)
			if err != nil {
				return errors.New(help.Header + err.Error())
			}
			return nil
		}
	}

	return fmt.Errorf("unknown arguments: %s", flag.MergeFlags(args))
}
