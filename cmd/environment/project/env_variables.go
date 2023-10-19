package project

import (
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	projPatterns "github.com/Red-Sock/rscli/plugins/project/patterns"
)

type envVariables struct {
	*env.Container
	envResources
}

const (
	userDefinedEnvVariablesComment       = "# user defined variables"
	globalUserDefinedEnvVariablesComment = "# global user defined variables"
)

func (e *envVariables) fetch(globalEnvFile *env.Container, pathToProjEnv string) error {
	dotEnvFilePath := path.Join(pathToProjEnv, projPatterns.ExampleFile+patterns.EnvFile.Name)
	envFile, err := os.ReadFile(dotEnvFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project .env example file "+dotEnvFilePath)
		}
	}

	e.Container, err = env.NewEnvContainer(nil)
	if err != nil {
		return errors.Wrap(err, "error initializing empty container")
	}

	example, err := env.NewEnvContainer(envFile)
	if err != nil {
		return errors.Wrap(err, "error creating project example .env file")
	}

	e.Container.AppendRaw(globalUserDefinedEnvVariablesComment, "")
	e.Container.Append(globalEnvFile.GetContent()...)

	e.Container.AppendRaw(userDefinedEnvVariablesComment, "")
	e.Container.Append(example.GetContent()...)

	e.envResources = newEnvManager(globalEnvFile)

	return nil
}

type envResources map[string]string

func newEnvManager(envContainer *env.Container) envResources {
	r := make(envResources)

	for _, item := range envContainer.GetContent() {
		switch {
		case strings.HasPrefix(item.Name, patterns.ResourceCapsPattern) &&
			len(item.Name) > len(patterns.ResourceCapsPattern):
			name := item.Name[len(patterns.ResourceCapsPattern)+1:]
			if name != "" {
				r[name] = item.Value
			}

		default:
			continue
		}
	}

	return r
}

func (e envResources) GetByName(envName string) string {
	for name, value := range e {
		if strings.HasSuffix(envName, name) {
			return value
		}
	}

	return ""
}
