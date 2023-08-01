package processor

import (
	"strconv"

	configstructs2 "github.com/Red-Sock/rscli/plugins/project/config/pkg/configstructs"
)

func DefaultRdsPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})
	if len(args) == 0 {
		args = append(args, "redis")
	}

	port := uint16(6379)

	for _, name := range args {
		out[name] = &configstructs2.Redis{
			Host: "0.0.0.0",
			Port: port,
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

	port := uint16(5432)

	for _, name := range args {
		out[name] = &configstructs2.Postgres{
			Name:    name,
			User:    name,
			Pwd:     name,
			Host:    "0.0.0.0",
			Port:    port,
			SSLMode: "disabled",
		}
		port++
	}

	return out
}

func DefaultHTTPPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})

	port := uint16(8080)
	name := "rest_api"

	if len(args) != 0 {
		for _, item := range args {
			p, err := strconv.ParseUint(item, 10, 16)
			if err == nil {
				port = uint16(p)
			}
		}
	}

	out[name] = &configstructs2.ServerOptions{
		Name:        name,
		Port:        port,
		CertPath:    "path/to/cert.crt",
		KeyPath:     "path/to/key.pem",
		ForceUseTLS: false,
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
	name := "grpc"
	out[name] = &configstructs2.ServerOptions{
		Name:        name,
		Port:        port,
		CertPath:    "path/to/cert.crt",
		KeyPath:     "path/to/key.pem",
		ForceUseTLS: false,
	}

	return out
}

func DefaultTelegramPattern(args []string) map[string]interface{} {
	out := make(map[string]interface{})

	tg := &configstructs2.Telegram{}
	out["tg"] = tg

	if len(args) != 0 {
		tg.ApiKey = args[0]
	}

	return out
}
