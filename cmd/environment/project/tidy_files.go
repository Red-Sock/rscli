package project

import (
	"bytes"
	"strconv"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/project"
)

func (e *Env) preTidyConfigFile() {
	projConfig, err := project.LoadProjectConfig(e.projPath, e.rscliConfig)
	if err != nil {
		return
	}

	for k := range e.Config.DataSources {
		if _, ok := projConfig.DataSources[k]; !ok {
			delete(e.Config.DataSources, k)
		}
	}

	for k, v := range projConfig.DataSources {
		if _, ok := e.Config.DataSources[k]; !ok {
			e.Config.DataSources[k] = v
		}
	}
}

func (e *Env) tidyServerAPIs(projName string, pm *ports.PortManager) error {
	opts, err := e.Config.GetServerOptions()
	if err != nil {
		return errors.Wrap(err, "error obtaining server options")
	}

	service, ok := e.Compose.Services[strings.ReplaceAll(projName, "-", "_")]
	if !ok {
		service = e.Compose.Services[patterns.ProjNamePattern]
	}

	for optName := range opts {
		portName := strings.ToUpper(projName) + "_" + strings.ToUpper(opts[optName].GetName()) + "_" + patterns.PortSuffix
		portName = strings.ReplaceAll(portName,
			"__", "_")

		externalPort := opts[optName].GetPort()

		if externalPort == 0 {
			continue
		}

		composePort := compose.AddEnvironmentBrackets(portName) + ":" + strconv.FormatUint(uint64(opts[optName].GetPort()), 10)
		portExists := false
		for _, item := range service.Ports {
			if item == composePort {
				portExists = true
				break
			}
		}
		if !portExists {
			service.Ports = append(service.Ports, composePort)
		}

		e.Environment.AppendRaw(portName, strconv.FormatUint(uint64(pm.GetNextPort(opts[optName].GetPort(), portName)), 10))
	}

	return nil
}

func (e *Env) preTidyEnvFile() {
	for _, envVar := range e.Environment.GetContent() {
		if envVar.Name == "" || envVar.Name[0] == '#' {
			e.Environment.Remove(envVar.Name)
		}
	}
}

func (e *Env) tidyMakeFile(projName string) {
	e.Makefile.Merge(e.globalMakefile)

	projNameCaps := strings.ToUpper(projName)

	{
		// tidy variables
		v := e.Makefile.GetVars().GetContent()

		for i := range v {
			v[i].Name = strings.ReplaceAll(v[i].Name, patterns.ProjNameCapsPattern, projNameCaps)
			switch v[i].Value {
			case patterns.AbsoluteProjectPathPattern:
				v[i].Value = strings.ReplaceAll(v[i].Value, patterns.AbsoluteProjectPathPattern, e.projPath)
			case patterns.PathToMain:
				v[i].Value = strings.ReplaceAll(v[i].Value, patterns.PathToMain, e.rscliConfig.Env.PathToMain)

			default:
				v[i].Value = renamer.ReplaceProjectNameStr(v[i].Value, projName)
			}
		}
	}

	{
		environments := make([]string, 0, len(e.Compose.Services))
		for name := range e.Compose.Services {
			if name != projName {
				environments = append(environments, name)
			}
		}

		rules := e.Makefile.GetRules()
		for i := range rules {
			if string(rules[i].Name) == patterns.MakefileEnvUpRuleName {
				envUpRule := e.globalMakefile.GetRuleByName(patterns.MakefileEnvUpRuleName)
				if envUpRule == nil {
					continue
				}

				if len(envUpRule.Commands) == 0 {
					continue
				}

				if len(rules[i].Commands) == 0 {
					rules[i].Commands = envUpRule.Commands
				}

				if !bytes.HasSuffix(envUpRule.Commands[0], []byte{' '}) {
					rules[i].Commands[0] = append(envUpRule.Commands[0], ' ')
				}

				rules[i].Commands[0] = append(rules[i].Commands[0], []byte(strings.Join(environments, " "))...)
			}
		}
	}
}
