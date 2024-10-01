package config

import (
	_ "embed"
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
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
	envPathToConfig = "RSCLI_PATH_TO_CONFIG"
	envPathToMain   = "RSCLI_PATH_TO_MAIN"

	envPathToProtoClients    = "RSCLI_PATH_TO_PROTO_CLIENTS"
	envPathToCompiledClients = "RSCLI_PATH_TO_COMPILED_CLIENTS"
	envPathToClients         = "RSCLI_PATH_TO_CLIENTS"

	envPathToServers          = "RSCLI_PATH_TO_SERVERS"
	envPathToServerDefinition = "RSCLI_PATH_TO_SERVER_DEFINITION"

	envPathToMigrations      = "RSCLI_PATH_TO_MIGRATIONS"
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
	PathToMain   string `yaml:"path_to_main"`
	PathToConfig string `yaml:"path_to_config"`

	PathsToCompiledClients []string `yaml:"paths_to_compiled_clients"`
	PathsToProtoClients    []string `yaml:"paths_to_proto_clients"`
	PathsToClients         []string `yaml:"paths_to_clients"`

	PathToServers          []string `yaml:"path_to_servers"`
	PathToServerDefinition string   `yaml:"path_to_server_definition"`

	PathToMigrations string `yaml:"path_to_migrations"`
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
		return errors.Wrap(err, "error obtaining config from custom file")
	}

	if configFromFile != nil {
		rsCliConfig = mergeConfigs(*configFromFile, rsCliConfig)
	}

	return nil
}

func getConfigFromEnvironment() (r RsCliConfig) {
	r.Env.PathToMain = os.Getenv(envPathToMain)
	r.Env.PathToConfig = os.Getenv(envPathToConfig)
	r.Env.PathToMigrations = os.Getenv(envPathToMigrations)
	r.Env.PathToServerDefinition = os.Getenv(envPathToServerDefinition)

	r.DefaultProjectGitPath = os.Getenv(envDefaultProjectGitPath)

	if pathToClients := strings.Split(os.Getenv(envPathToClients), ","); pathToClients[0] != "" {
		r.Env.PathsToClients = pathToClients
	}

	if pathToServers := strings.Split(os.Getenv(envPathToServers), ","); pathToServers[0] != "" {
		r.Env.PathToServers = pathToServers
	}

	if pathToProtoClients := strings.Split(os.Getenv(envPathToProtoClients), ","); pathToProtoClients[0] != "" {
		r.Env.PathsToProtoClients = pathToProtoClients
	}

	if pathToCompiledClients := strings.Split(os.Getenv(envPathToCompiledClients), ","); pathToCompiledClients[0] != "" {
		r.Env.PathsToCompiledClients = pathToCompiledClients
	}
	return
}

func getConfigFromFile(cmd *cobra.Command) (*RsCliConfig, error) {
	if cmd == nil {
		return nil, nil
	}

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

	if master.Env.PathToConfig == "" {
		master.Env.PathToConfig = slave.Env.PathToConfig
	}

	if master.DefaultProjectGitPath == "" {
		master.DefaultProjectGitPath = slave.DefaultProjectGitPath
	}

	if len(master.Env.PathsToClients) == 0 {
		master.Env.PathsToClients = slave.Env.PathsToClients
	}

	if len(master.Env.PathToServers) == 0 {
		master.Env.PathToServers = slave.Env.PathToServers
	}

	if master.Env.PathToMigrations == "" {
		master.Env.PathToMigrations = slave.Env.PathToMigrations
	}

	if master.Env.PathToServerDefinition == "" {
		master.Env.PathToServerDefinition = slave.Env.PathToServerDefinition
	}

	if len(master.Env.PathsToProtoClients) == 0 {
		master.Env.PathsToProtoClients = slave.Env.PathsToProtoClients
	}

	if len(master.Env.PathsToCompiledClients) == 0 {
		master.Env.PathsToCompiledClients = slave.Env.PathsToCompiledClients
	}

	return master
}
