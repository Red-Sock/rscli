package project

import (
	"bytes"
	"strings"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

func (e *ProjEnv) tidyMakeFile(projName string) {
	e.Makefile.Merge(e.globalMakefile)

	projNameCaps := strings.ToUpper(projName)

	{
		// tidy variables
		v := e.Makefile.GetVars().GetContent()

		for i := range v {
			v[i].Name = strings.ReplaceAll(v[i].Name, patterns.ProjNameCapsPattern, projNameCaps)
			switch v[i].Value {
			case patterns.AbsoluteProjectPathPattern:
				v[i].Value = strings.ReplaceAll(v[i].Value, patterns.AbsoluteProjectPathPattern, e.srcProjPath)
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
