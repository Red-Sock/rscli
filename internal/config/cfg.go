package config

import (
	_ "embed"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/flag"
)

const configFilename = "rscli.yaml"

const customPathToConfig = "-cfg"

//go:embed rscli.yaml
var example []byte

type RsCliConfig struct {
	Env struct {
		PathToMain   string `yaml:"path_to_main"`
		PathToConfig string `yaml:"path_to_config"`
	} `yaml:"env"`
}

func ReadConfig(args []string) (*RsCliConfig, error) {
	flags := flag.ParseArgs(args)

	cfgFile, err := flag.ExtractOneValueFromFlags(flags, customPathToConfig)
	if err != nil {
		return nil, err
	}
	var file []byte
	if cfgFile != "" {
		file, err = os.ReadFile(cfgFile)
		if err != nil {
			return nil, err
		}
	} else {
		file = example
	}

	var c RsCliConfig

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
