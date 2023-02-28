package config

import (
	"github.com/Red-Sock/rscli/pkg/flag"
	"gopkg.in/yaml.v3"
	"os"

	_ "embed"
)

const configFilename = "rscli.yaml"

const customPathToConfig = "-cfg"

//go:embed rscli.yaml
var example []byte

type Config struct {
	Env struct {
		PathToMain    string `yaml:"path_to_main"`
		PathToClients string `yaml:"path_to_clients"`
	} `yaml:"env"`
}

func ReadConfig(args []string) (*Config, error) {
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

	var c Config

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
