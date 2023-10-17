package env

import (
	"strings"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
)

type envManager struct {
	env       *env.Container
	resources resources
}

func newEnvManager(envContainer *env.Container) *envManager {
	em := &envManager{
		env:       envContainer,
		resources: make(resources),
	}

	for _, item := range envContainer.GetContent() {
		switch {
		case strings.HasPrefix(item.Name, patterns.ResourceCapsPattern) &&
			len(item.Name) > len(patterns.ResourceCapsPattern):
			name := item.Name[len(patterns.ResourceCapsPattern)+1:]
			if name != "" {
				em.resources[name] = item.Value
			}

		default:
			continue
		}
	}

	return em
}

type resources map[string]string

func (e resources) GetByName(envName string) string {
	for name, value := range e {
		if strings.HasSuffix(envName, name) {
			return value
		}
	}

	return ""
}
