package project

import (
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
)

func (e *Env) tidyResources(pm *ports.PortManager, projName string) error {
	tr := tidyResources{
		config:                      e.Config,
		compose:                     e.Compose,
		environment:                 e.Environment,
		composePatterns:             e.ComposePatterns,
		environmentResourcePatterns: e.environmentResourcePatterns,
		pm:                          pm,
		projName:                    projName,
	}
	return tr.tidyResources()
}

type tidyResources struct {
	config      *config.Config
	compose     *compose.Compose
	environment *env.Container

	composePatterns             compose.PatternManager
	environmentResourcePatterns envResourcePattern
	pm                          *ports.PortManager

	projName string
}

func (e *tidyResources) tidyResources() error {
	dependencies, err := e.composePatterns.GetServiceDependencies(e.config)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+e.config.AppInfo.Name)
	}

	for _, resource := range dependencies {

		err = e.tidyResource(e.projName, resource)
		if err != nil {
			return errors.Wrap(err, "error tiding resource "+resource.GetName())
		}

		e.compose.AppendService(resource.GetName(), resource.GetCompose())
	}

	return nil
}

func (e *tidyResources) tidyResource(projName string, resource compose.Pattern) (err error) {
	patternEnv := resource.GetEnvs().Content()
	envMap := make(map[string]string, len(patternEnv))
	for idx := range patternEnv {
		newEnvName := e.getResourceName(patternEnv[idx].Name, resource.GetName(), projName)

		basicEnvName, newEnvValue := e.getDefaultValue(patternEnv[idx].Name, resource.GetType())

		if strings.HasSuffix(newEnvName, patterns.PortSuffix) {
			newEnvValue = e.getPort(newEnvName, patternEnv[idx].Value)
		} else {
			newEnvValue = renamer.ReplaceProjectNameStr(newEnvValue, projName)
		}

		envMap[basicEnvName] = newEnvValue

		resource.RenameVariable(patternEnv[idx].Name, newEnvName)

		e.environment.AppendRaw(newEnvName, newEnvValue)
	}

	e.config.DataSources[resource.GetName()] = e.tidyConfig(resource, envMap)

	return nil
}

func (e *tidyResources) tidyConfig(resource compose.Pattern, env map[string]string) interface{} {
	switch resource.GetType() {
	case resources.DataSourcePostgres:
		pgConf := resources.Postgres{}
		pgConf.FillFromEnv(env)
		return pgConf
	default:
		return nil
	}
}

func (e *tidyResources) getResourceName(varName, resName, projName string) string {
	newEnvName := strings.ReplaceAll(varName,
		patterns.ResourceNameCapsPattern, strings.ToUpper(resName))

	newEnvName = strings.ReplaceAll(newEnvName,
		"__", "_")
	return renamer.ReplaceProjectNameStr(newEnvName, projName)
}

func (e *tidyResources) getPort(envName, envVal string) string {
	port, err := strconv.ParseUint(envVal, 10, 16)
	if err != nil {
		port = 10_000
	}

	return strconv.FormatUint(uint64(e.pm.GetNextPort(uint16(port), envName)), 10)
}

func (e *tidyResources) getDefaultValue(resName, resType string) (basicEnvName, envValue string) {
	basicEnvName = strings.ReplaceAll(resName,
		patterns.ResourceNameCapsPattern, strings.ToUpper(resType))
	basicEnvName = strings.ReplaceAll(basicEnvName,
		patterns.ProjNameCapsPattern+"_", "")

	return basicEnvName, e.environmentResourcePatterns.GetByName(basicEnvName)
}
