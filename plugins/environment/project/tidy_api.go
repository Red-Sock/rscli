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
		p, ok := e.globalComposePatternManager.Patterns[envpatterns.ProjNamePattern]
		if !ok {
			return errors.New("no pattern for service was found")
		}

		service = &p.ContainerDefinition
	}

	for i := range opts {
		portName := strings.ToUpper(e.projName) + "_" + strings.ToUpper(opts[i].GetName()) + "_" + envpatterns.PortSuffix
		portName = strings.ReplaceAll(portName,
			"__", "_")

		internalPort := opts[i].GetPort()

		if internalPort == 0 {
			continue
		}

		composePort := compose.AddEnvironmentBrackets(portName) + ":" + strconv.FormatUint(uint64(internalPort), 10)
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

		e.Environment.AppendRaw(portName, strconv.FormatUint(uint64(e.globalPortManager.GetNextPort(opts[i].GetPort(), portName)), 10))
	}

	e.Compose.AppendService(e.projName, *service)

	return nil
}
