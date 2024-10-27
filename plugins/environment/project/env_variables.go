package project

import (
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	projPatterns "github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
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
	dotEnvFilePath := path.Join(pathToProjEnv, projPatterns.ExampleFile+envpatterns.EnvFile.Name)
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

	globalValidEnvs := make([]env.Variable, 0, len(globalEnvFile.Content)/2)

	for _, item := range globalEnvFile.GetContent() {
		if len(item.Name) == 0 {
			continue
		}

		if strings.HasPrefix(item.Name, "###") {
			continue
		}

		if strings.HasPrefix(item.Name, envpatterns.ResourceCapsPattern) {
			continue
		}

		globalValidEnvs = append(globalValidEnvs, item)
	}

	if len(globalValidEnvs) != 0 {
		e.Container.AppendRaw(globalUserDefinedEnvVariablesComment, "")
		e.Container.Append(globalValidEnvs...)
	}

	if len(example.Content) != 0 {
		e.Container.AppendRaw(userDefinedEnvVariablesComment, "")
		e.Container.Append(example.GetContent()...)
	}

	e.envResources = newEnvManager(globalEnvFile)

	return nil
}

type envResources map[string]string

func newEnvManager(envContainer *env.Container) envResources {
	r := make(envResources)

	for _, item := range envContainer.GetContent() {
		switch {
		case strings.HasPrefix(item.Name, envpatterns.ResourceCapsPattern) &&
			len(item.Name) > len(envpatterns.ResourceCapsPattern):
			name := item.Name[len(envpatterns.ResourceCapsPattern)+1:]
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
