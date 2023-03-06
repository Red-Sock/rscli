package actions

import (
	"bytes"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
	"sort"
	"strings"

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

	keysFromCfg := strings.Split(string(keys), "\n")
	sort.Slice(keysFromCfg, func(i, j int) bool {
		return keysFromCfg[i] < keysFromCfg[j]
	})
	keys = []byte(strings.Join(keysFromCfg, "\n"))
	cfgKeysFile := make([]byte, len(patterns.ConfigKeys))
	copy(cfgKeysFile, patterns.ConfigKeys)
	body := append(cfgKeysFile[:bytes.Index(cfgKeysFile, []byte("// _start_of_consts_to_replace"))], keys...)
	body = append(body, cfgKeysFile[bytes.Index(cfgKeysFile, []byte("// _end_of_consts_to_replace")):]...)
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
