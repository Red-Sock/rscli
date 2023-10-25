package project

import (
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/utils/copier"
)

var ErrNoProjectComposePattern = errors.New("")

func (e *ProjEnv) tidyService() error {
	srcPattern, ok := e.globalComposePatternManager.Patterns[envpatterns.ProjNamePattern]
	if !ok {
		return ErrNoProjectComposePattern
	}

	var pattern compose.Pattern
	err := copier.Copy(srcPattern, &pattern)
	if err != nil {
		return errors.Wrap(err, "error coping proj pattern")
	}

	srvs, err := e.Config.GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error getting server options")
	}

	for _, s := range srvs {
		port := s.GetPort()

		pattern.ContainerDefinition.Ports = append(
			pattern.ContainerDefinition.Ports,
			fmt.Sprintf("%d:%d", e.globalPortManager.GetNextPort(port, e.projName+"_"+s.GetName()), port),
		)
	}

	e.Compose.AppendService(e.projName, pattern.ContainerDefinition)

	return nil
}
