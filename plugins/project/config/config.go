package config

import (
	"os"
	"strings"
	"time"

	"github.com/Red-Sock/trace-errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
	"github.com/Red-Sock/rscli/plugins/project/config/server"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type Config struct {
	AppInfo     AppInfo                `yaml:"app_info"`
	Server      map[string]interface{} `yaml:"server,omitempty"`
	DataSources map[string]interface{} `yaml:"data_sources,omitempty"`

	pth string `yaml:"-"`
}

type AppInfo struct {
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	StartupDuration time.Duration `yaml:"startup_duration"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]interface{}{},
		DataSources: map[string]interface{}{},
	}
}

type ConnectionOptions struct {
	Type string
	Name string

	ConnectionString string
}

func ReadConfig(pth string) (*Config, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}

	c := &Config{
		pth: pth,
	}
	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding config to struct")
	}

	if c.DataSources == nil {
		c.DataSources = make(map[string]interface{})
	}

	if c.Server == nil {
		c.Server = map[string]interface{}{}
	}

	return c, nil
}

func (c *Config) BuildTo(cfgFile string) error {
	f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error opening cfg file")
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(*c)
	if err != nil {
		return errors.Wrap(err, "error encoding config to "+cfgFile)
	}

	return nil
}

// GetDataSourceFolders extracts data sources folders
func (c *Config) GetDataSourceFolders() (*folder.Folder, error) {
	out := &folder.Folder{
		Name: patterns.ClientsFolder,
	}

	datasourceUniqueTypes := map[string]struct{}{}

	for dataSourceName := range c.DataSources {
		dataSourceType := strings.Split(dataSourceName, "_")[0]

		connFolder, err := patterns.GetDatasourceClientFile(dataSourceType)
		if err != nil {
			return nil, errors.Wrap(err, "error obtaining client conn files for datasource")
		}

		if _, ok := datasourceUniqueTypes[dataSourceType]; ok {
			continue
		}

		datasourceUniqueTypes[dataSourceType] = struct{}{}

		out.Inner = append(out.Inner, connFolder)
	}

	return out, nil
}

func (c *Config) GetServerFolders() ([]*folder.Folder, error) {
	out := make([]*folder.Folder, 0, len(c.Server))

	for serverName := range c.Server {
		serverPattern, err := patterns.GetServerFiles(strings.Split(serverName, "_")[0])
		if err != nil {
			return nil, errors.Wrap(err, "error getting server files")
		}

		if serverPattern.Validators != nil {
			serverPattern.Validators(&serverPattern.F, serverName)
		}

		serverPattern.F.Name = serverName

		out = append(out, &serverPattern.F)
	}
	return out, nil
}

func (c *Config) GetDataSourceOptions() ([]resources.Resource, error) {
	if len(c.DataSources) == 0 {
		return nil, nil
	}
	confResources := make([]resources.Resource, 0, len(c.DataSources))
	for dsn, data := range c.DataSources {

		rsrs, err := resources.ParseResource(dsn, data)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing resource")
		}

		confResources = append(confResources, rsrs)
	}

	return confResources, nil
}

func (c *Config) GetServerOptions() ([]server.Server, error) {
	out := make([]server.Server, 0, len(c.Server))

	for name, content := range c.Server {
		opt, err := server.ParseServerOption(name, content)
		if err != nil {
			return out, errors.Wrap(err, "error parsing server option")
		}

		out = append(out, opt)
	}

	return out, nil
}

func (c *Config) GenerateGoConfigKeys(prefix string) ([]byte, error) {
	envKeys := ParseKeysFromEnv(prefix)

	keysFromCfg, err := c.keysFromConfig()
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

func (c *Config) ExtractName() string {
	return c.AppInfo.Name
}

func (c *Config) GetProjInfo() AppInfo {
	return c.AppInfo
}

func (c *Config) keysFromConfig() (map[string]string, error) {
	cfgBytes, err := yaml.Marshal(*c)
	if err != nil {
		return nil, err
	}

	cfg := make(cfgKeysBuilder)
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, err
	}

	vars, err := cfg.extractVariables("", cfg)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]string, len(cfg))
	for _, v := range vars {
		variables[cases.SnakeToCamel(v)] = v
	}

	return variables, nil
}

func (c *Config) GetPath() string {
	return c.pth
}
