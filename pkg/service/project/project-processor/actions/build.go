package actions

import (
	"bytes"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/interfaces"
	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/patterns"

	"github.com/Red-Sock/rscli/pkg/folder"
)

func BuildConfigGoFolder(p interfaces.Project) error {
	out := []*folder.Folder{
		{
			Name: "config.go",
			Content: []byte(
				strings.ReplaceAll(patterns.Configurator, "{{projectNAME_}}", strings.ToUpper(p.GetName())),
			),
		},
	}

	keys, err := p.GetConfig().GenerateGoConfigKeys(p.GetName())
	if err != nil {
		return err
	}
	cfgKeysStr := string(patterns.ConfigKeys)

	body := append(patterns.ConfigKeys[:strings.Index(cfgKeysStr, "// _start_of_consts_to_replace")], keys...)
	body = append(body, patterns.ConfigKeys[strings.Index(cfgKeysStr, "// _end_of_consts_to_replace"):]...)
	body = bytes.ReplaceAll(body, []byte("// _end_of_consts_to_replace"), []byte(""))

	if len(keys) != 0 {
		out = append(out,
			&folder.Folder{
				Name:    "keys.go",
				Content: body,
			})
	}

	p.GetFolder().AddWithPath(
		[]string{
			"internal",
			"config",
		},
		out...,
	)

	return nil
}

func BuildProject(p interfaces.Project) error {

	ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build("")
	if err != nil {
		return err
	}
	return nil
}
