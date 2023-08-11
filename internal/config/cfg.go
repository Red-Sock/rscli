package config

import (
	_ "embed"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/flag"
)

const (
	configFilename          = "rscli.yaml"
	customPathToConfig      = "--rscli-cfg"
	environmentPathToConfig = "RSCLI_CONFIG_PATH"
)

const (
	envPathToConfig          = "RSCLI_PATH_TO_CONFIG"
	envPathToMain            = "RSCLI_PATH_TO_MAIN"
	envDefaultProjectGitPath = "RSCLI_DEFAULT_PROJECT_GIT_PATH"
)

//go:embed rscli.yaml
var builtInConfig []byte

type RsCliConfig struct {
	Env                   Environment `yaml:"env"`
	DefaultProjectGitPath string      `yaml:"default_project_git_path"`
}

type Environment struct {
	PathToMain   string `yaml:"path_to_main"`
	PathToConfig string `yaml:"path_to_config"`
}

func GetConfig() *RsCliConfig {
	var builtInConf RsCliConfig
	err := yaml.Unmarshal(builtInConfig, &builtInConf)
	if err != nil {
		panic(errors.Wrap(err, "error parsing built in config file. This is serious issue and MUST BE fixed ASAP\n\n\n Rocky (: "))
	}

	out := mergeConfigs(getConfigFromEnvironment(), builtInConf)

	externalConf, err := getConfigFromFile()
	if externalConf != nil {
		out = mergeConfigs(*externalConf, out)
	}

	return &out
}

func getConfigFromEnvironment() (r RsCliConfig) {
	r.Env.PathToMain = os.Getenv(envPathToMain)
	r.Env.PathToConfig = os.Getenv(envPathToConfig)

	r.DefaultProjectGitPath = os.Getenv(envDefaultProjectGitPath)

	return
}

func getConfigFromFile() (*RsCliConfig, error) {
	cfgFilePath, _ := flag.ExtractOneValueFromFlags(flag.ParseArgs(os.Args[1:]), customPathToConfig)
	if cfgFilePath == "" {
		cfgFilePath = os.Getenv(environmentPathToConfig)
	}

	if cfgFilePath == "" {
		exePath, _ := os.Executable()
		cfgFilePath = path.Join(path.Dir(exePath), configFilename)
	}

	file, err := os.ReadFile(cfgFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errors.Wrap(err, "error reading file from FS")
	}

	if len(file) == 0 {
		return nil, nil
	}

	var externalConf RsCliConfig
	err = yaml.Unmarshal(file, &externalConf)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling config from: "+cfgFilePath)
	}

	return &externalConf, nil
}

func mergeConfigs(master, slave RsCliConfig) RsCliConfig {
	if master.Env.PathToMain == "" {
		master.Env.PathToMain = slave.Env.PathToMain
	}

	if master.Env.PathToConfig == "" {
		master.Env.PathToConfig = slave.Env.PathToConfig
	}

	if master.DefaultProjectGitPath == "" {
		master.DefaultProjectGitPath = slave.DefaultProjectGitPath
	}

	return master
}
