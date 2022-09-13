package config

import (
	"github.com/Red-Sock/rscli/pkg/config"
	"strconv"
)

func DefaultRdsPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})
	if len(args) == 0 {
		args = append(args, "redis")
	}

	port := 6379

	for _, name := range args {
		out[name] = &config.Redis{
			Host: "0.0.0.0",
			Port: strconv.Itoa(port),
		}
		port++
	}

	return out
}

func DefaultPgPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})
	if len(args) == 0 {
		args = append(args, "postgres")
	}

	port := 5432

	for _, name := range args {
		out[name] = &config.Postgres{
			Name: name,
			User: name,
			Pwd:  name,
			Host: "0.0.0.0",
			Port: strconv.Itoa(port),
		}
		port++
	}

	return out
}
