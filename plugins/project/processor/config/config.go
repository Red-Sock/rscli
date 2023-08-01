package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/helpers/copier"
	"github.com/Red-Sock/rscli/pkg/folder"
	configstructs2 "github.com/Red-Sock/rscli/plugins/project/config/pkg/configstructs"
	"github.com/Red-Sock/rscli/plugins/project/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

const apiInfoKey = _const.AppKey + "_info"

type ProjectConfig struct {
	Path string

	Values map[string]interface{}

	appInfo *configstructs2.AppInfo
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

func (c *ProjectConfig) Rebuild(p interfaces.Project) error {
	{
		// App info
		bts, err := json.Marshal(c.appInfo)
		if err != nil {
			return err
		}

		trg := map[string]interface{}{}

		err = json.Unmarshal(bts, &trg)
		if err != nil {
			return err
		}

		c.Values[apiInfoKey] = trg
	}

	separatedPath := strings.Split(c.GetPath(), string(filepath.Separator))

	cfgFile := p.GetFolder().GetByPath(patterns.ConfigsFolder, separatedPath[len(separatedPath)-1])
	if cfgFile == nil {
		return errors.New("cannot find config at " + c.GetPath())
	}

	var err error
	cfgFile.Content, err = yaml.Marshal(c.Values)
	if err != nil {
		return err
	}

	return nil
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

// GetDataSourceFolders extracts data sources information from _consts file and parses it as folders in project
func (c *ProjectConfig) GetDataSourceFolders() (*folder.Folder, error) {
	// extract connections from data source part
	connectionSettings, ok := c.Values[_const.DataSourceKey]
	if !ok {
		return nil, nil
	}

	var ds map[string]interface{}
	ds, ok = connectionSettings.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	out := &folder.Folder{
		Name: patterns.ClientsFolder,
	}

	connTypeAdded := map[string]struct{}{}

	for dsn := range ds {
		dataSourceType := strings.Split(dsn, "_")[0]
		connFile := patterns.DatasourceClients[dataSourceType]
		if connFile == nil {
			return nil, errors.New(fmt.Sprintf("unknown data source %s. "+
				"DataSource should start with name of source (e.g redis, postgres)"+
				"and (or) be followed by \"_\" symbol if needed (e.g redis_shard1, postgres_replica2)", dsn))
		}

		if _, ok := connTypeAdded[dataSourceType]; ok {
			continue
		}

		connTypeAdded[dataSourceType] = struct{}{}

		out.Inner = append(out.Inner, &folder.Folder{Name: dsn, Inner: connFile})
	}

	// extract connections from server part
	connectionSettings, ok = c.Values[_const.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	ds, ok = connectionSettings.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	for dsn := range ds {
		connFile, ok := patterns.DatasourceClients[dsn]
		if !ok {
			continue
		}

		if _, ok := connTypeAdded[dsn]; ok {
			continue
		}

		connTypeAdded[dsn] = struct{}{}

		out.Inner = append(out.Inner, &folder.Folder{Name: dsn, Inner: connFile})
	}

	return out, nil
}

func (c *ProjectConfig) GetServerFolders() ([]*folder.Folder, error) {
	serverOpts, ok := c.Values[_const.ServerOptsKey]
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

		if ptrn.Validators != nil {
			ptrn.Validators(&serverF, serverName)
		}

		serverF.Name = serverName

		out = append(out, &serverF)
	}
	return out, nil
}

func (c *ProjectConfig) GetProjInfo() (*configstructs2.AppInfo, error) {
	if c.appInfo != nil {
		return c.appInfo, nil
	}

	appInfoMap, ok := c.Values[apiInfoKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = appInfoMap.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	bts, err := yaml.Marshal(so)
	if err != nil {
		return nil, err
	}

	c.appInfo = &configstructs2.AppInfo{}

	err = yaml.Unmarshal(bts, c.appInfo)
	if err != nil {
		return nil, err
	}

	return c.appInfo, nil
}

func (c *ProjectConfig) GetServerOptions() ([]configstructs2.ServerOptions, error) {
	serverOpts, ok := c.Values[_const.ServerOptsKey]
	if !ok {
		return nil, nil
	}

	var so map[string]interface{}
	so, ok = serverOpts.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	out := make([]configstructs2.ServerOptions, 0, len(so))

	for key, item := range so {
		servOpt := configstructs2.ServerOptions{
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

func (c *ProjectConfig) GetDataSourceOptions() (out []configstructs2.ConnectionOptions, err error) {
	dataSources, ok := c.Values[_const.DataSourceKey]
	if !ok {
		return nil, nil
	}

	var ds map[string]interface{}
	ds, ok = dataSources.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	for dsn, data := range ds {
		dataSourceType := strings.Split(dsn, "_")[0]
		switch dataSourceType {
		case _const.SourceNamePostgres:
			var pgDSN configstructs2.Postgres
			err = copier.Copy(data, &pgDSN)
			if err != nil {
				return nil, err
			}

			out = append(out, configstructs2.ConnectionOptions{
				Type:             _const.SourceNamePostgres,
				Name:             dsn,
				ConnectionString: fmt.Sprintf(_const.PostgresConnectionString, pgDSN.User, pgDSN.Pwd, pgDSN.Host, pgDSN.Port, pgDSN.Name),
			})
		}
	}

	return out, nil
}

func (c *ProjectConfig) ExtractName() (string, error) {
	var cfg configstructs2.Config
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
