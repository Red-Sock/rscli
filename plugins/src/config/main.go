package main

import (
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/internal/utils/input"
	"github.com/Red-Sock/rscli/pkg/service/project/config-processor/config"
	"github.com/pkg/errors"
	"os"
	"strings"
)

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return "config"
}

func (p *plugin) Run(args []string) error {
	c, err := config.Run(args)
	if err != nil {
		return errors.Wrap(err, "error Running config ")
	}

	err = c.TryWrite()
	if err != nil {
		if err != os.ErrExist {
			return errors.Wrap(err, "error Writing config")
		}

		answ := input.GetAnswer("file " + c.GetPath() + " already exists. Do you want to override it? Y(es)/N(o)")
		answ = strings.ToLower(strings.Replace(answ, "\n", "", -1))

		if !strings.HasPrefix(answ, "y") {
			println("aborting config creation. " + randomizer.GoodGoodBuy())
		}

		err = c.ForceWrite()
		if err != nil {
			return errors.Wrap(err, "error forcing writing")

		}
	}
	println("successfully create config at " + c.GetPath() + ". " + randomizer.GoodGoodBuy())
	return nil
}
