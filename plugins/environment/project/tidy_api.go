package project

import (
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/envpatterns"
)

func (e *ProjEnv) tidyServerAPIs() error {
	opts, err := e.Config.GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining server options")
	}

	service, ok := e.Compose.Services[strings.ReplaceAll(e.projName, "-", "_")]
	if !ok {
		service = e.Compose.Services[envpatterns.ProjNamePattern]
	}

	for optName := range opts {
		portName := strings.ToUpper(e.projName) + "_" + strings.ToUpper(opts[optName].GetName()) + "_" + envpatterns.PortSuffix
		portName = strings.ReplaceAll(portName,
			"__", "_")

		externalPort := opts[optName].GetPort()

		if externalPort == 0 {
			continue
		}

		composePort := compose.AddEnvironmentBrackets(portName) + ":" + strconv.FormatUint(uint64(opts[optName].GetPort()), 10)
		portExists := false
		for _, item := range service.Ports {
			if item == composePort {
				portExists = true
				break
			}
		}
		if !portExists {
			service.Ports = append(service.Ports, composePort)
		}

		e.Environment.AppendRaw(portName, strconv.FormatUint(uint64(e.globalPortManager.GetNextPort(opts[optName].GetPort(), portName)), 10))
	}

	return nil
}
