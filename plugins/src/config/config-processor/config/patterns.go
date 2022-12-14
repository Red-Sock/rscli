package config

import (
	"strconv"

	"github.com/Red-Sock/rscli/pkg/config"
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

func DefaultHTTPPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})

	port := uint16(8080)

	if len(args) != 0 {
		p, err := strconv.ParseUint(args[0], 10, 16)
		if err == nil {
			port = uint16(p)
		}
	}

	out["rest_api"] = &config.RestApi{
		Port: strconv.FormatUint(uint64(port), 10),
	}

	return out
}

func DefaultGRPCPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})

	port := uint16(50051)

	if len(args) != 0 {
		p, err := strconv.ParseUint(args[0], 10, 16)
		if err == nil {
			port = uint16(p)
		}
	}

	out["grpc"] = &config.RestApi{
		Port: strconv.FormatUint(uint64(port), 10),
	}

	return out
}
