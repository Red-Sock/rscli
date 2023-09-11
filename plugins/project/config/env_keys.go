package config

import (
	"os"
	"strings"
)

// ParseKeysFromEnv extracts keys from env
func ParseKeysFromEnv(prefix string) map[string]string {
	envToGo := func(in string) (out string) {
		keyWords := strings.Split(in, "_")
		for idx := range keyWords {
			keyWords[idx] = strings.ToUpper(keyWords[idx][:1]) + keyWords[idx][1:]
		}
		return strings.Join(keyWords, "_")
	}

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
		values[envToGo(name)] = name
	}

	return values
}
