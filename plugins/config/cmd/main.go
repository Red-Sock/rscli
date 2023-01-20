package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/plugins/config/processor"
)

const PluginName = "config"

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return PluginName
}

func (p *plugin) Run(args map[string][]string) error {
	c, err := processor.Run(args)
	if err != nil {
		return errors.Wrap(err, "error Running config ")
	}

	err = c.TryWrite()
	if err != nil {
		if err != os.ErrExist {
			return errors.Wrap(err, "error Writing config")
		}

		println("file " + c.GetPath() + " already exists. Do you want to override it? Y(es)/N(o)")
		res, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		answ := strings.ToLower(strings.Replace(string(res), "\n", "", -1))

		if !strings.HasPrefix(answ, "y") {
			println("aborting config creation. ")
		}

		err = c.ForceWrite()
		if err != nil {
			return errors.Wrap(err, "error forcing writing")

		}
	}
	println("successfully create config at " + c.GetPath() + ". ")
	return nil
}
