package config

import (
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func KeysFromConfig(pathToConfig string) (map[string]string, error) {
	cfgBytes, err := os.ReadFile(pathToConfig)
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
		parts := strings.Split(v[1:], "_")
		for i := range parts {
			if slices.Contains(initialisms, parts[i]) {
				parts[i] = strings.ToUpper(parts[i])
			} else {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}
		variables[strings.Join(parts, "")] = v[1:]
	}

	return variables, nil
}

type cfgKeysBuilder map[string]interface{}

func (c *cfgKeysBuilder) extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(cfgKeysBuilder); ok {
			values, err := c.extractVariables(prefix+"_"+k, newMap)
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

var initialisms = []string{"acl", "api", "ascii", "cpu", "css", "dns",
	"eof", "guid", "html", "http", "https", "id",
	"ip", "json", "qps", "ram", "rpc", "sla",
	"smtp", "sql", "ssh", "tcp", "tls", "ttl",
	"udp", "ui", "gid", "uid", "uuid", "uri",
	"url", "utf8", "vm", "xml", "xmpp", "xsrf",
	"xss", "sip", "rtp", "amqp", "db", "ts"}
