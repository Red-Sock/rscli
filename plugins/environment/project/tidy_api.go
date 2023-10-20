package project

import (
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/environment/project/compose"
	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"
)

func (e *ProjEnv) tidyServerAPIs(projName string) error {
	opts, err := e.Config.GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining server options")
	}

	service, ok := e.Compose.Services[strings.ReplaceAll(projName, "-", "_")]
	if !ok {
		service = e.Compose.Services[envpatterns.ProjNamePattern]
	}

	for optName := range opts {
		portName := strings.ToUpper(projName) + "_" + strings.ToUpper(opts[optName].GetName()) + "_" + envpatterns.PortSuffix
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
