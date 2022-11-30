package project

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/config"
	"gopkg.in/yaml.v3"
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
		name: "data",
	}

	for dsn := range ds {
		out.inner = append(out.inner, &Folder{
			name: dsn,
		})
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
func generateConfigKeys(prefix string) []byte {
	envKeys := configKeysFromEnv(prefix)
	configKeys := convertConfigKeysToGoConstName()

	for e, v := range envKeys {
		if _, ok := configKeys[e]; !ok {
			configKeys[e] = v
		}
	}
	sb := &strings.Builder{}
	for key, v := range configKeys {
		sb.WriteString(key + `="` + v + `"`)
	}

	return []byte(sb.String())
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
func convertConfigKeysToGoConstName() map[string]string {
	return nil
}

// generates names for go const env keys
func convertEnvVarToGoConstName(in string) (out string) {
	keyWords := strings.Split(in, "_")
	for idx := range keyWords {
		keyWords[idx] = strings.ToUpper(keyWords[idx][:1]) + keyWords[idx][1:]
	}
	return strings.Join(keyWords, "_")
}
