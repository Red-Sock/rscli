package project

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/service/config"
	"gopkg.in/yaml.v3"
)

type dataSourcePrefix string

const (
	RedisDataSourcePrefix    dataSourcePrefix = "redis"
	PostgresDataSourcePrefix dataSourcePrefix = "postgres"
)

type serverOptsPrefix string

const (
	RESTServerPrefix serverOptsPrefix = "rest"
	GRPCServerPrefix serverOptsPrefix = "grpc"
)

type Config struct {
	path string

	values map[string]interface{}
}

// NewProjectConfig - constructor for configuration of project
func NewProjectConfig(p string) *Config {
	return &Config{
		path: p,
	}
}

// prepares self to be worked on
func (c *Config) parseSelf() error {
	c.values = make(map[string]interface{})

	bytes, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &c.values)
	if err != nil {
		return err
	}
	return nil
}

// extracts data sources information from config file and parses it as folders in project
func (c *Config) extractDataSources() (*Folder, error) {
	dataSources, ok := c.values[config.DataSourceKey]
	if !ok {
		return nil, nil
	}

	var ds map[string]interface{}
	ds, ok = dataSources.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := &Folder{
		name: "clients",
	}

	for dsn := range ds {
		file := datasourceClients[dataSourcePrefix(strings.Split(dsn, "_")[0])]
		if file == nil {
			return nil, errors.New(fmt.Sprintf("unknown data source %s. "+
				"DataSource should start with name of source (e.g redis, postgres)"+
				"and (or) be followed by \"_\" symbol if needed (e.g redis_shard1, postgres_replica2)", dsn))
		}

		out.inner = append(out.inner, &Folder{
			name: dsn,
			inner: []*Folder{
				{
					name:    "conn.go",
					content: file,
				},
			},
		})
	}

	return out, nil
}

func (c *Config) extractServerOptions() (*Folder, error) {
	serverOpts, ok := c.values[config.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = serverOpts.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := &Folder{
		name: "transport",
	}

	for serverName := range so {
		files := serverOptsPatterns[serverOptsPrefix(strings.Split(serverName, "_")[0])]
		if files == nil {
			return nil, errors.New(fmt.Sprintf("unknown server option %s. "+
				"Server Option should start with type of server (e.g rest, grpc)"+
				"and (or) be followed by \"_\" symbol if needed (e.g rest_v1, grpc_proxy)", serverName))
		}
		if len(files) == 0 {
			continue
		}
		serverFolder := &Folder{
			name: serverName,
		}
		for name, content := range files {
			serverFolder.inner = append(serverFolder.inner, &Folder{
				name:    name,
				content: content,
			})
		}

		out.inner = append(out.inner, serverFolder)
	}
	return out, nil
}

// tries to find path to configuration in same directory
func findConfigPath() (pth string, err error) {
	currentDir := "./"

	var dirs []os.DirEntry
	dirs, err = os.ReadDir(currentDir)
	if err != nil {
		return "", err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = path.Join(currentDir, config.DefaultDir)
			break
		}
	}

	if pth == "" {
		return "", nil
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return "", err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, config.FileName) {
			pth = path.Join(pth, name)
			break
		}
	}

	return pth, nil
}

// extracts one value from number of flags
func extractOneValueFromFlags(flagsArgs map[string][]string, flags ...string) (string, error) {
	var name []string
	for _, f := range flags {
		var ok bool
		name, ok = flagsArgs[f]
		if ok {
			break
		}
	}

	if name == nil {
		return "", nil
	}

	if len(name) == 0 {
		return "", fmt.Errorf("%w expected 1 got 0 ", ErrNoArgumentsSpecifiedForFlag)
	}

	if len(name) > 1 {
		return "", fmt.Errorf("%w expected 1 got %d", ErrFlagHasTooManyArguments, len(name))
	}

	return name[0], nil
}

// generates file with constant config keys
func generateConfigKeys(prefix, pathToConfig string) ([]byte, error) {
	envKeys := configKeysFromEnv(prefix)
	configKeys, err := convertConfigKeysToGoConstName(pathToConfig)
	if err != nil {
		return nil, err
	}
	for e, v := range envKeys {
		if _, ok := configKeys[e]; !ok {
			configKeys[e] = v
		}
	}
	sb := &strings.Builder{}

	for key, v := range configKeys {
		sb.WriteString(key + "= \"" + v + "\"\n")
	}
	
	return []byte(sb.String()), nil
}

// extracts keys from env
func configKeysFromEnv(prefix string) map[string]string {
	listEnv := os.Environ()
	envs := make([]string, 0, 10)
	for _, e := range listEnv {
		if strings.HasPrefix(e, prefix+"_") {
			envs = append(envs, e[len(prefix)+1:])
		}
	}

	values := map[string]string{}

	for _, e := range envs {
		nAv := strings.Split(e, "=")
		if len(nAv) != 2 {
			continue
		}
		name := nAv[0]
		values[convertEnvVarToGoConstName(name)] = name
	}

	return values
}

// generates names for go const config keys
func convertConfigKeysToGoConstName(pathToConfig string) (map[string]string, error) {
	cfgBytes, err := os.ReadFile(pathToConfig)
	if err != nil {
		return nil, err
	}
	cfg := make(map[string]interface{})
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, err
	}

	vars, err := extractVariables("", cfg)
	if err != nil {
		return nil, err
	}
	variables := make(map[string]string, len(cfg))
	for _, v := range vars {
		parts := strings.Split(v[1:], "_")
		for i := range parts {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
		variables[strings.Join(parts, "")] = v[1:]
	}

	return variables, nil
}

func extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(map[string]interface{}); ok {
			values, err := extractVariables(prefix+"_"+k, newMap)
			if err != nil {
				return nil, err
			}
			out = append(out, values...)
		} else {
			out = append(out, prefix+"_"+k)
		}
	}
	return out, nil
}

// generates names for go const env keys
func convertEnvVarToGoConstName(in string) (out string) {
	keyWords := strings.Split(in, "_")
	for idx := range keyWords {
		keyWords[idx] = strings.ToUpper(keyWords[idx][:1]) + keyWords[idx][1:]
	}
	return strings.Join(keyWords, "_")
}
