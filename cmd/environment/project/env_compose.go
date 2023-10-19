package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

type envCompose struct {
	*compose.Compose
}

func (e *envCompose) fetch(pathToProjectEnv string) error {
	projectEnvComposeFilePath := path.Join(pathToProjectEnv, patterns.DockerComposeFile.Name)
	composeFile, err := os.ReadFile(projectEnvComposeFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project env docker-compose file "+projectEnvComposeFilePath)
		}
	}

	if len(composeFile) == 0 {
		globalEnvComposeFilePath := path.Join(path.Dir(pathToProjectEnv), patterns.DockerComposeFile.Name)
		composeFile, err = os.ReadFile(globalEnvComposeFilePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global docker-compose file "+globalEnvComposeFilePath)
			}
		}
	}

	if len(composeFile) == 0 {
		projName := path.Base(pathToProjectEnv)
		composeFile = renamer.ReplaceProjectName(patterns.DockerComposeFile.Content, projName)
	}

	e.Compose, err = compose.NewComposeAssembler(composeFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}
