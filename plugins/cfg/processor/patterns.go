package processor

import (
	"github.com/Red-Sock/rscli/plugins/cfg/pkg/structs"
	"strconv"
)

func DefaultRdsPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})
	if len(args) == 0 {
		args = append(args, "redis")
	}

	port := 6379

	for _, name := range args {
		out[name] = &structs.Redis{
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
		out[name] = &structs.Postgres{
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

	out["rest_api"] = &structs.RestApi{
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

	out["grpc"] = &structs.RestApi{
		Port: strconv.FormatUint(uint64(port), 10),
	}

	return out
}
