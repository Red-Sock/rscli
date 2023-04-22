package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	config "github.com/Red-Sock/rscli/plugins/config/pkg/const"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/config/pkg/structs"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type ProjectConfig struct {
	Path string

	Values map[string]interface{}
}

// NewProjectConfig - constructor for configuration of project
func NewProjectConfig(p string) (*ProjectConfig, error) {
	c := &ProjectConfig{
		Path: p,
	}
	err := c.ParseSelf()

	return c, err
}

func (c *ProjectConfig) SetPath(pth string) {
	c.Path = pth
}

func (c *ProjectConfig) GetPath() string {
	return c.Path
}

// ParseSelf prepares self to be worked on
func (c *ProjectConfig) ParseSelf() error {
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
func (c *ProjectConfig) GetDataSourceFolders() (*folder.Folder, error) {
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
		file := patterns.DatasourceClients[strings.Split(dsn, "_")[0]]
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

func (c *ProjectConfig) GetServerFolders() ([]*folder.Folder, error) {
	serverOpts, ok := c.Values[config.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = serverOpts.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := make([]*folder.Folder, 0, len(so))

	for serverName := range so {
		ptrn, ok := patterns.ServerOptsPatterns[strings.Split(serverName, "_")[0]]
		if !ok {
			return nil, errors.New(fmt.Sprintf("unknown server option %s. "+
				"Server Option should start with type of server (e.g rest, grpc)"+
				"and (or) be followed by \"_\" symbol if needed (e.g rest_v1, grpc_proxy)", serverName))
		}

		// it must be a file of folder with files|folders
		if len(ptrn.F.Inner)+len(ptrn.F.Content) == 0 {
			continue
		}

		// copy pattern
		var serverF folder.Folder
		mb, err := json.Marshal(ptrn.F)
		if err != nil {
			return nil, errors.Wrap(err, "error marshalling folder pattern")
		}
		err = json.Unmarshal(mb, &serverF)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling copy of folder pattern")
		}

		ptrn.Validators(&serverF, serverName)
		serverF.Name = serverName

		out = append(out, &serverF)
	}
	return out, nil
}

type ServerOptions struct {
	Name string
	Port string `json:"port"`
}

func (c *ProjectConfig) GetServerOptions() ([]ServerOptions, error) {
	serverOpts, ok := c.Values[config.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = serverOpts.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	out := make([]ServerOptions, 0, len(so))

	for key, item := range so {
		servOpt := ServerOptions{
			Name: key,
		}
		bts, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bts, &servOpt)
		if err != nil {
			return nil, err
		}

		out = append(out, servOpt)
	}

	return out, nil
}

func (c *ProjectConfig) ExtractName() (string, error) {
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

func (c *ProjectConfig) GenerateGoConfigKeys(prefix string) ([]byte, error) {
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
		sb.WriteString(key + " = \"" + v + "\"\n")
	}

	return []byte(sb.String()), nil
}

func (c *ProjectConfig) GetTemplate() ([]byte, error) {
	b, err := yaml.Marshal(c.Values)
	if err != nil {
		return nil, err
	}

	return b, nil
}
