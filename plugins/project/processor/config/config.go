package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/folder"
	structs "github.com/Red-Sock/rscli/plugins/config/pkg/structs"
	config "github.com/Red-Sock/rscli/plugins/config/processor"
	"github.com/Red-Sock/rscli/plugins/project/processor/consts"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type Config struct {
	Path string

	Values map[string]interface{}
}

// NewProjectConfig - constructor for configuration of project
func NewProjectConfig(p string) (*Config, error) {
	c := &Config{
		Path: p,
	}
	err := c.ParseSelf()

	return c, err
}

func (c *Config) SetPath(pth string) {
	c.Path = pth
}

func (c *Config) GetPath() string {
	return c.Path
}

// ParseSelf prepares self to be worked on
func (c *Config) ParseSelf() error {
	if c.Values != nil {
		return nil
	}
	c.Values = make(map[string]interface{})

	bytes, err := os.ReadFile(c.Path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &c.Values)
	if err != nil {
		return err
	}
	return nil
}

// ExtractDataSources extracts data sources information from config file and parses it as folders in project
func (c *Config) ExtractDataSources() (*folder.Folder, error) {
	dataSources, ok := c.Values[config.DataSourceKey]
	if !ok {
		return nil, nil
	}

	var ds map[string]interface{}
	ds, ok = dataSources.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := &folder.Folder{
		Name: "clients",
	}

	for dsn := range ds {
		file := patterns.DatasourceClients[consts.DataSourcePrefix(strings.Split(dsn, "_")[0])]
		if file == nil {
			return nil, errors.New(fmt.Sprintf("unknown data source %s. "+
				"DataSource should start with name of source (e.g redis, postgres)"+
				"and (or) be followed by \"_\" symbol if needed (e.g redis_shard1, postgres_replica2)", dsn))
		}

		out.Inner = append(out.Inner, &folder.Folder{
			Name: dsn,
			Inner: []*folder.Folder{
				{
					Name:    "conn.go",
					Content: file,
				},
			},
		})
	}

	return out, nil
}

func (c *Config) ExtractServerOptions() (*folder.Folder, error) {
	serverOpts, ok := c.Values[config.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = serverOpts.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := &folder.Folder{
		Name: "transport",
	}

	for serverName := range so {
		files := patterns.ServerOptsPatterns[consts.ServerOptsPrefix(strings.Split(serverName, "_")[0])]
		if files == nil {
			return nil, errors.New(fmt.Sprintf("unknown server option %s. "+
				"Server Option should start with type of server (e.g rest, grpc)"+
				"and (or) be followed by \"_\" symbol if needed (e.g rest_v1, grpc_proxy)", serverName))
		}
		if len(files) == 0 {
			continue
		}
		serverFolder := &folder.Folder{
			Name: serverName,
		}
		for name, content := range files {
			serverFolder.Inner = append(serverFolder.Inner, &folder.Folder{
				Name:    name,
				Content: content,
			})
		}

		out.Inner = append(out.Inner, serverFolder)
	}
	return out, nil
}

func (c *Config) ExtractName() (string, error) {
	var cfg structs.Config
	bytes, err := yaml.Marshal(c.Values)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return "", err
	}

	return cfg.AppInfo.Name, nil
}

func (c *Config) GenerateGoConfigKeys(prefix string) ([]byte, error) {
	envKeys := KeysFromEnv(prefix)

	keysFromCfg, err := KeysFromConfig(c.Path)
	if err != nil {
		return nil, err
	}

	for e, v := range envKeys {
		if _, ok := keysFromCfg[e]; !ok {
			keysFromCfg[e] = v
		}
	}
	sb := &strings.Builder{}

	for key, v := range keysFromCfg {
		sb.WriteString(key + "= \"" + v + "\"\n")
	}

	return []byte(sb.String()), nil
}
