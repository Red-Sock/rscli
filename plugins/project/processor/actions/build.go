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
	cfgKeysFile := make([]byte, 0, len(patterns.ConfigKeys))

	cfgKeysFile = append(cfgKeysFile, patterns.ConfigKeys[:bytes.Index(patterns.ConfigKeys, []byte("// _start_of_consts_to_replace"))]...)
	cfgKeysFile = append(cfgKeysFile, keys...)
	endOfConstsToReplaceBytes := []byte("// _end_of_consts_to_replace")
	cfgKeysFile = append(cfgKeysFile, patterns.ConfigKeys[bytes.Index(patterns.ConfigKeys, endOfConstsToReplaceBytes)+len(endOfConstsToReplaceBytes):]...)

	if len(keys) != 0 {
		out = append(out,
			&folder.Folder{
				Name:    "keys.go",
				Content: cfgKeysFile,
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
