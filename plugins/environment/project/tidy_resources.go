package project

import (
	"sort"
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/environment/project/compose"
	"github.com/Red-Sock/rscli/plugins/environment/project/compose/env"
	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
)

func (e *ProjEnv) tidyResources(projName string, enableService bool) error {
	dependencies, err := e.globalComposePatternManager.GetServiceDependencies(e.Config.Config)
	if err != nil {
		return errors.Wrap(err, "error getting dependencies for service "+e.Config.Config.AppInfo.Name)
	}

	envs := make([]env.Container, 0, len(dependencies))

	for _, resource := range dependencies {
		var resourceEnv env.Container
		resourceEnv, err = e.tidyResource(projName, resource, enableService)
		if err != nil {
			return errors.Wrap(err, "error tiding resource "+resource.GetName())
		}

		e.Compose.AppendService(resource.GetName(), resource.GetCompose())

		envs = append(envs, resourceEnv)
	}

	sort.Slice(envs, func(i, j int) bool {
		if len(envs[i].Content) == 0 {
			return true
		}
		if len(envs[j].Content) == 0 {
			return false
		}

		return envs[i].Content[0].Value > envs[j].Content[0].Name
	})

	for _, item := range envs {
		e.Environment.Append(item.Content...)
	}

	for name := range e.Compose.Services {
		foundInConfig := false

		for _, cfgRes := range dependencies {
			if name == cfgRes.GetName() || name == projName {
				foundInConfig = true
				break
			}
		}

		if !foundInConfig {
			delete(e.Compose.Services, name)
		}
	}

	return nil
}

func (e *ProjEnv) tidyResource(projName string, resource compose.Pattern, enableService bool) (container env.Container, err error) {
	patternEnv := resource.GetEnvs().GetContent()
	envMap := make(map[string]string, len(patternEnv))

	container.AppendRaw("# "+resource.GetName(), "")

	for idx := range patternEnv {
		var ev env.Variable
		ev.Name = e.getResourceName(patternEnv[idx].Name, resource.GetName(), projName)

		var basicEnvName string
		basicEnvName, ev.Value = e.getDefaultValue(patternEnv[idx].Name, resource.GetType())

		if strings.HasSuffix(ev.Name, envpatterns.PortSuffix) {
			ev.Value = e.getPort(ev.Name, patternEnv[idx].Value)
		} else {
			ev.Value = renamer.ReplaceProjectNameStr(ev.Value, projName)
		}

		envMap[basicEnvName] = ev.Value

		resource.RenameVariable(patternEnv[idx].Name, ev.Name)

		container.Append(ev)
	}

	hostName := strings.ToUpper(projName+"_"+resource.GetName()) + envpatterns.HostEnvSuffix
	hostValue := ""
	if enableService {
		hostValue = resource.GetName()
	} else {
		hostValue = envpatterns.Localhost
	}

	envMap[strings.ToUpper(resource.GetType()+envpatterns.HostEnvSuffix)] = hostValue
	container.AppendRaw(hostName, hostValue)

	e.Config.DataSources[resource.GetName()] = e.tidyResourceConfig(resource, envMap)

	return container, nil
}

func (e *ProjEnv) tidyResourceConfig(resource compose.Pattern, env map[string]string) interface{} {
	switch resource.GetType() {
	case resources.DataSourcePostgres:
		pgConf := resources.Postgres{}
		_ = pgConf.FillFromEnv(env)
		return pgConf
	case resources.DataSourceRedis:
		rdsConf := resources.Redis{}
		_ = rdsConf.FillFromEnv(env)
		return rdsConf
	default:
		return nil
	}
}

func (e *ProjEnv) getResourceName(varName, resName, projName string) string {
	newEnvName := strings.ReplaceAll(varName,
		envpatterns.ResourceNameCapsPattern, strings.ToUpper(resName))

	newEnvName = strings.ReplaceAll(newEnvName,
		"__", "_")
	return renamer.ReplaceProjectNameStr(newEnvName, projName)
}

func (e *ProjEnv) getPort(envName, envVal string) string {
	port, err := strconv.ParseUint(envVal, 10, 16)
	if err != nil {
		port = 10_000
	}

	return strconv.FormatUint(uint64(e.globalPortManager.GetNextPort(uint16(port), envName)), 10)
}

func (e *ProjEnv) getDefaultValue(resName, resType string) (basicEnvName, envValue string) {
	basicEnvName = strings.ReplaceAll(resName,
		envpatterns.ResourceNameCapsPattern, strings.ToUpper(resType))
	basicEnvName = strings.ReplaceAll(basicEnvName,
		envpatterns.ProjNameCapsPattern+"_", "")

	return basicEnvName, e.Environment.envResources.GetByName(basicEnvName)
}
