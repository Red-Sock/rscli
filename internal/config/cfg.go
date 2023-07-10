package config

import (
	_ "embed"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/flag"
)

const (
	configFilename          = "rscli.yaml"
	customPathToConfig      = "--rscli-cfg"
	environmentPathToConfig = "RSCLI_CONFIG_PATH"
)

const (
	pathToConfig = "RSCLI_PATH_TO_CONFIG"
	pathToMain   = "RSCLI_PATH_TO_MAIN"
)

//go:embed rscli.yaml
var builtInConfig []byte

type RsCliConfig struct {
	Env struct {
		PathToMain   string `yaml:"path_to_main"`
		PathToConfig string `yaml:"path_to_config"`
	} `yaml:"env"`
}

func GetConfig() (*RsCliConfig, error) {
	var builtInConf RsCliConfig
	err := yaml.Unmarshal(builtInConfig, &builtInConf)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config file")
	}

	envConf := getConfigFromEnvironment()

	out := mergeConfigs(envConf, builtInConf)

	cfgFilePath := getConfigPathFromArgs()
	if cfgFilePath == "" {
		cfgFilePath = getConfigPathFromExecutable()
	}

	if cfgFilePath == "" {
		cfgFilePath = getConfigPathFromEnvironment()
	}

	file, err := os.ReadFile(cfgFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errors.Wrap(err, "error reading file from FS")
	}

	if len(file) == 0 {
		return &out, nil
	}

	var externalConf RsCliConfig
	err = yaml.Unmarshal(file, &externalConf)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling config from: "+cfgFilePath)
	}

	out = mergeConfigs(externalConf, out)

	return &out, nil
}

func getConfigFromEnvironment() (r RsCliConfig) {
	r.Env.PathToMain = os.Getenv(pathToMain)
	r.Env.PathToConfig = os.Getenv(pathToConfig)

	return
}

func getConfigPathFromArgs() string {
	cfgFilePath, _ := flag.ExtractOneValueFromFlags(flag.ParseArgs(os.Args[1:]), customPathToConfig)
	return cfgFilePath
}

func getConfigPathFromExecutable() string {
	exePath, _ := os.Executable()
	return path.Join(path.Dir(exePath), configFilename)
}

func getConfigPathFromEnvironment() string {
	return os.Getenv(environmentPathToConfig)
}

func mergeConfigs(master, slave RsCliConfig) RsCliConfig {
	if master.Env.PathToMain == "" {
		master.Env.PathToMain = slave.Env.PathToMain
	}

	if master.Env.PathToConfig == "" {
		master.Env.PathToConfig = slave.Env.PathToConfig
	}

	return master
}
