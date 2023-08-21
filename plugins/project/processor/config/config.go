package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/helpers/copier"
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type Config struct {
	AppInfo     AppInfo                  `yaml:"app_info"`
	Server      map[string]ServerOptions `yaml:"server,omitempty"`
	DataSources map[string]interface{}   `yaml:"data_sources,omitempty"`
}

type AppInfo struct {
	Name            string        `yaml:"name"`
	Version         string        `yaml:"version"`
	StartupDuration time.Duration `yaml:"startup_duration"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Server:      map[string]ServerOptions{},
		DataSources: map[string]interface{}{},
	}
}

type ServerOptions struct {
	Name        string
	Port        uint16 `yaml:"port"`
	CertPath    string `yaml:"cert_path"`
	KeyPath     string `yaml:"key_path"`
	ForceUseTLS bool   `yaml:"force_use_tls"`
}

type ConnectionOptions struct {
	Type string
	Name string

	ConnectionString string
}

func ParseConfig(pth string) (*Config, error) {
	panic("TODO")
}

func (c *Config) BuildTo(cfgFile string) error {

	f, err := os.Open(cfgFile)
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
	// extract connections from data source part

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

func (c *Config) GetProjInfo() AppInfo {
	return c.AppInfo
}

func (c *Config) GetServerOptions() map[string]ServerOptions {
	return c.Server
}

func (c *Config) GetDataSourceOptions() (out []ConnectionOptions, err error) {
	for dsn, data := range c.DataSources {
		dataSourceType := strings.Split(dsn, "_")[0]
		switch dataSourceType {
		case patterns.SourceNamePostgres:
			var pgDSN Postgres
			err = copier.Copy(data, &pgDSN)
			if err != nil {
				return nil, err
			}

			out = append(out, ConnectionOptions{
				Type:             patterns.SourceNamePostgres,
				Name:             dsn,
				ConnectionString: fmt.Sprintf(PostgresConnectionString, pgDSN.User, pgDSN.Pwd, pgDSN.Host, pgDSN.Port, pgDSN.Name),
			})
		default:
			return nil, errors.Wrapf(err, "unknown datasource type %s", dataSourceType)
		}
	}

	return out, nil
}

func (c *Config) ExtractName() string {
	return c.AppInfo.Name
}

func (c *Config) GenerateGoConfigKeys(prefix string) ([]byte, error) {
	envKeys := KeysFromEnv(prefix)

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

func (c *Config) GetTemplate() ([]byte, error) {
	b, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	return b, nil
}
