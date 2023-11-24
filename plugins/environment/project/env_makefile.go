package project

import (
	"bytes"
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/makefile"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

type envMakefile struct {
	*makefile.Makefile
}

func (e *envMakefile) fetch(envProjPath string) (err error) {
	e.Makefile, err = makefile.ReadMakeFile(path.Join(envProjPath, envpatterns.Makefile.Name))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error getting makefile")
		}

		e.Makefile = makefile.MewEmptyMakefile()
	}

	return nil
}

func (e *ProjEnv) tidyMakeFile() {
	e.Makefile.Merge(e.globalMakefile)

	projNameCaps := strings.ToUpper(e.projName)

	{
		// tidy variables
		v := e.Makefile.GetVars().GetContent()

		for i := range v {
			v[i].Name = strings.ReplaceAll(v[i].Name, envpatterns.ProjNameCapsPattern, projNameCaps)
			switch v[i].Value {
			case envpatterns.AbsoluteProjectPathPattern:
				v[i].Value = strings.ReplaceAll(v[i].Value, envpatterns.AbsoluteProjectPathPattern, e.pathToProjSrc)
			case envpatterns.PathToMain:
				v[i].Value = strings.ReplaceAll(v[i].Value, envpatterns.PathToMain, e.rscliConfig.Env.PathToMain)

			default:
				v[i].Value = renamer.ReplaceProjectNameStr(v[i].Value, e.projName)
			}
		}
	}

	{
		environments := make([]string, 0, len(e.Compose.Services))
		for name := range e.Compose.Services {
			if name != e.projName {
				environments = append(environments, name)
			}
		}

		rules := e.Makefile.GetRules()
		for i := range rules {
			if string(rules[i].Name) == envpatterns.MakefileEnvUpRuleName {
				envUpRule := e.globalMakefile.GetRuleByName(envpatterns.MakefileEnvUpRuleName)
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
