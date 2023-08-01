package go_actions

import (
	"bytes"
	"sort"
	"strings"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

type BuildGoConfigFolderAction struct {
}

func (a BuildGoConfigFolderAction) Do(p interfaces.Project) error {
	out := []*folder.Folder{
		{
			Name: patterns.ConfigsFolder,
			Content: []byte(
				strings.ReplaceAll(patterns.ConfiguratorFile, "{{projectNAME_}}", strings.ToUpper(p.GetShortName())),
			),
		},
	}

	keys, err := p.GetConfig().GenerateGoConfigKeys(p.GetName())
	if err != nil {
		return err
	}

	keysFromCfg := strings.Split(string(keys), "\n")
	keysFromCfg = keysFromCfg[:len(keysFromCfg)-1]
	sort.Slice(keysFromCfg, func(i, j int) bool {
		return keysFromCfg[i] < keysFromCfg[j]
	})
	keys = []byte(strings.Join(keysFromCfg, "\n\t"))
	cfgKeysFile := make([]byte, 0, len(patterns.ConfigKeysFile))

	cfgKeysFile = append(cfgKeysFile, patterns.ConfigKeysFile[:bytes.Index(patterns.ConfigKeysFile, []byte("// _start_of_consts_to_replace"))]...)
	cfgKeysFile = append(cfgKeysFile, keys...)
	endOfConstsToReplaceBytes := []byte("// _end_of_consts_to_replace")
	cfgKeysFile = append(cfgKeysFile, patterns.ConfigKeysFile[bytes.Index(patterns.ConfigKeysFile, endOfConstsToReplaceBytes)+len(endOfConstsToReplaceBytes):]...)

	if len(keys) != 0 {
		out = append(out,
			&folder.Folder{
				Name:    "keys.go",
				Content: cfgKeysFile,
			})
	}

	p.GetFolder().ForceAddWithPath(
		[]string{
			"internal",
			"config",
		},
		out...,
	)

	return nil
}
func (a BuildGoConfigFolderAction) NameInAction() string {
	return "Building config folder"
}

type BuildProjectAction struct {
}

func (a BuildProjectAction) Do(p interfaces.Project) error {

	ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build()
	if err != nil {
		return err
	}
	return nil
}
func (a BuildProjectAction) NameInAction() string {
	return "Building project"
}
