package project

import (
	"fmt"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"os"
	"path"
	"strings"
)

func (p *Project) tryFindConfig(args map[string][]string) error {
	pth, err := extractCfgPathFromFlags(args)
	if err != nil {
		return err
	}

	if pth != "" {
		p.CfgPath = pth
		return nil
	}

	currentDir := "./"
	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		return err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = path.Join(currentDir, config.DefaultDir)
			break
		}
	}

	if pth == "" {
		return ErrNoConfigNoAppNameFlag
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, config.FileName) {
			pth = path.Join(pth, name)
			break
		}
	}
	p.CfgPath = pth
	return nil
}

func extractDataSources(ds map[string]interface{}) (folder, error) {
	out := folder{
		name: "data",
	}

	for dsn := range ds {
		out.inner = append(out.inner, folder{
			name: dsn,
		})
	}

	return out, nil
}

func extractCfgPathFromFlags(flagsArgs map[string][]string) (string, error) {
	name, ok := flagsArgs[FlagCfgPath]
	if !ok {
		name, ok = flagsArgs[FlagCfgPathShort]
		if !ok {
			return "", nil
		}
	}
	if len(name) == 0 {
		return "", fmt.Errorf("%w expected 1 got 0 ", ErrNoArgumentsSpecifiedForFlag)
	}

	if len(name) > 1 {
		return "", fmt.Errorf("%w expected 1 got %d", ErrFlagHasTooManyArguments, len(name))
	}

	return name[0], nil
}

func generateConfig(projectName string) []folder {

	out := []folder{
		{
			name: "config.go",
			content: []byte(
				strings.ReplaceAll(configurator, "{{projectNAME}}", strings.ToUpper(projectName)),
			),
		},
	}

	keys := generateConfigKeys(projectName)
	if len(keys) != 0 {
		out = append(out,
			folder{
				name:    "keys.go",
				content: keys,
			})
	}

	return out
}

func readEnvironment(prefix string) []string {
	listEnv := os.Environ()
	values := make([]string, 0, 10)
	for _, e := range listEnv {
		if strings.HasPrefix(e, prefix+"_") {
			values = append(values, e[len(prefix)+1:])
		}
	}
	return values
}

func generateConfigKeys(prefix string) []byte {
	envKeys := configKeysFromEnv(prefix)
	configKeys := configKeysFromConfig()

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

func configKeysFromEnv(prefix string) map[string]string {
	envVals := readEnvironment(prefix)

	values := map[string]string{}

	for _, e := range envVals {
		nAv := strings.Split(e, "=")
		if len(nAv) != 2 {
			continue
		}
		name := nAv[0]
		values[convertEnvVarToGoConstName(name)] = name
	}

	return values
}

func configKeysFromConfig() map[string]string {
	return nil
}

func convertEnvVarToGoConstName(in string) (out string) {
	keyWords := strings.Split(in, "_")
	for idx := range keyWords {
		keyWords[idx] = strings.ToUpper(keyWords[idx][:1]) + keyWords[idx][1:]
	}
	return strings.Join(keyWords, "_")
}
