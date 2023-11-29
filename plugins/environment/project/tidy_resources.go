package project

import (
	"sort"
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

func (e *ProjEnv) tidyResources(enableService bool) error {
	sort.Slice(e.Config.AppConfig.Resources, func(i, j int) bool {
		return e.Config.AppConfig.Resources[i].GetName() > e.Config.AppConfig.Resources[j].GetName()
	})

	for idx := range e.Config.AppConfig.Resources {
		err := e.tidyResource(e.projName, idx, enableService)
		if err != nil {
			return errors.Wrap(err, "error tiding resource "+
				e.Config.AppConfig.Resources[idx].GetName())
		}
	}

	for name := range e.Compose.Services {
		foundInConfig := false

		for _, cfgRes := range e.Config.AppConfig.Resources {
			if name == cfgRes.GetName() || name == e.projName {
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

func (e *ProjEnv) tidyResource(projName string, resourceIdx int, enableService bool) error {
	resource, err := e.globalComposePatternManager.GetServiceDependencies(e.Config.Resources[resourceIdx])
	if err != nil {
		return errors.Wrapf(err, "error getting resource pattern for type: %s, with name %s",
			e.Config.Resources[resourceIdx].GetType(),
			e.Config.Resources[resourceIdx].GetName())
	}

	if resource == nil {
		return nil
	}

	patternEnv := resource.GetEnvs().GetContent()

	envMap := make(map[string]string, len(patternEnv))

	envVars := make([]env.Variable, 0, len(patternEnv)+1)
	envVars = append(envVars, env.Variable{Name: "# " + resource.GetName()})

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
		envVars = append(envVars, ev)
	}

	hostName := strings.ToUpper(projName+"_"+resource.GetName()) + envpatterns.HostEnvSuffix
	hostValue := ""
	if enableService {
		hostValue = resource.GetName()
	} else {
		hostValue = envpatterns.Localhost
	}

	envMap[strings.ToUpper(resource.GetType()+envpatterns.HostEnvSuffix)] = hostValue

	envVars = append(envVars, env.Variable{Name: hostName, Value: hostValue})

	e.Environment.Append(envVars...)
	e.Compose.AppendService(resource.GetName(), resource.GetCompose())

	err = e.tidyResourceConfig(resourceIdx, resource, envMap)
	if err != nil {
		return errors.Wrap(err, "error tidy resource config")
	}

	return nil
}

func (e *ProjEnv) tidyResourceConfig(resourceIdx int, resource *compose.Pattern, env map[string]string) error {
	res := resources.GetResourceByName(resource.GetType())

	err := res.FromEnv(env)
	if err != nil {
		return errors.Wrap(err, "error filling config from env")
	}
	e.Config.Resources[resourceIdx] = res

	return nil
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
