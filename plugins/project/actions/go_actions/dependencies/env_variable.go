package dependencies

import (
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
)

type EnvVariable struct {
	dep dependencyBase
}

func envVariable(dep dependencyBase) Dependency {
	return &EnvVariable{
		dep: dep,
	}
}

func (e *EnvVariable) AppendToProject(proj Project) error {
	name := "new_variable"

	for _, env := range proj.GetConfig().Environment {
		if env.Value.Value() == name {
			return nil
		}
	}

	proj.GetConfig().Environment = append(proj.GetConfig().Environment,
		environment.MustNewVariable(name, "new string value"))

	return nil
}
