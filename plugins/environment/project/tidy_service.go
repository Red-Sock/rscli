package project

import (
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/utils/copier"
)

var ErrNoProjectComposePattern = rerrors.New("")

func (e *ProjEnv) tidyService() error {
	srcPattern, ok := e.globalComposePatternManager.Patterns[envpatterns.ProjNamePattern]
	if !ok {
		return ErrNoProjectComposePattern
	}

	var pattern compose.Pattern
	err := copier.Copy(srcPattern, &pattern)
	if err != nil {
		return rerrors.Wrap(err, "error coping proj pattern")
	}

	// TODO
	//for _, s := range e.Config.Servers {
	//port := s.GetPort()

	//pattern.ContainerDefinition.Ports = append(
	//	pattern.ContainerDefinition.Ports,
	//	fmt.Sprintf("%d:%d", e.globalPortManager.GetNextPort(port, e.projName+"_"+s.GetName()), port),
	//)
	//}

	e.Compose.AppendService(e.projName, pattern.ContainerDefinition)

	return nil
}
