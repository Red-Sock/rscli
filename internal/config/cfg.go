package config

import (
	_ "embed"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/flag"
)

const configFilename = "rscli.yaml"

const customPathToConfig = "--rscli-cfg"

//go:embed rscli.yaml
var example []byte

type RsCliConfig struct {
	Env struct {
		PathToMain   string `yaml:"path_to_main"`
		PathToConfig string `yaml:"path_to_config"`
	} `yaml:"env"`
}

func GetConfig() (*RsCliConfig, error) {
	flags := flag.ParseArgs(os.Args[1:])

	cfgFilePath, err := flag.ExtractOneValueFromFlags(flags, customPathToConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting config value from arguments")
	}

	var file []byte
	if cfgFilePath != "" {
		file, err = os.ReadFile(cfgFilePath)
		if err != nil {
			return nil, errors.Wrap(err, "error reading file from FS")
		}
	} else {
		file = example
	}

	var c RsCliConfig

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config file")
	}

	return &c, nil
}
