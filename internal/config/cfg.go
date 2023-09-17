package config

import (
	_ "embed"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	CustomPathToConfig = "rscli-cfg"
)
const (
	configFilename          = "rscli.yaml"
	environmentPathToConfig = "RSCLI_CONFIG_PATH"
)

const (
	envPathToConfig          = "RSCLI_PATH_TO_CONFIG"
	envPathToMain            = "RSCLI_PATH_TO_MAIN"
	envPathToClients         = "RSCLI_PATH_TO_CLIENTS"
	envDefaultProjectGitPath = "RSCLI_DEFAULT_PROJECT_GIT_PATH"
)

//go:embed rscli.yaml
var builtInConfig []byte

type RsCliConfig struct {
	Env                   Project `yaml:"env"`
	DefaultProjectGitPath string  `yaml:"default_project_git_path"`
}

var rsCliConfig RsCliConfig

type Project struct {
	PathToMain         string   `yaml:"path_to_main"`
	PathToConfigFolder string   `yaml:"path_to_config"`
	PathsToClients     []string `yaml:"paths_to_clients"`
}

func GetConfig() *RsCliConfig {
	return &rsCliConfig
}

func InitConfig(cmd *cobra.Command, _ []string) error {
	err := yaml.Unmarshal(builtInConfig, &rsCliConfig)
	if err != nil {
		panic(errors.Wrap(err, "error parsing built in config file. This is serious issue and MUST BE fixed A$A₽\n\n\n Like Rocky\n\n\n\n in a way (: "))
	}

	rsCliConfig = mergeConfigs(getConfigFromEnvironment(), rsCliConfig)

	configFromFile, err := getConfigFromFile(cmd)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if configFromFile != nil {
		rsCliConfig = mergeConfigs(*configFromFile, rsCliConfig)
	}

	return nil
}

func getConfigFromEnvironment() (r RsCliConfig) {
	r.Env.PathToMain = os.Getenv(envPathToMain)
	r.Env.PathToConfigFolder = os.Getenv(envPathToConfig)

	r.DefaultProjectGitPath = os.Getenv(envDefaultProjectGitPath)
	pathToClients := strings.Split(os.Getenv(envPathToClients), ",")
	if pathToClients[0] != "" {
		r.Env.PathsToClients = pathToClients
	}
	return
}

func getConfigFromFile(cmd *cobra.Command) (*RsCliConfig, error) {
	cfgFilePath := cmd.Flag(CustomPathToConfig).Value.String()

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

	if master.Env.PathToConfigFolder == "" {
		master.Env.PathToConfigFolder = slave.Env.PathToConfigFolder
	}

	if master.DefaultProjectGitPath == "" {
		master.DefaultProjectGitPath = slave.DefaultProjectGitPath
	}

	if len(master.Env.PathsToClients) == 0 {
		master.Env.PathsToClients = slave.Env.PathsToClients
	}

	return master
}
