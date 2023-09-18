package go_actions

import (
	"bytes"
	"path"
	"sort"
	"strings"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p interfaces.Project) error {
	out := []*folder.Folder{
		{
			Name: path.Join(patterns.InternalFolder, patterns.ConfigsFolder, patterns.ConfigFileName),
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
				Name:    path.Join(patterns.InternalFolder, patterns.ConfigsFolder, patterns.ConfigKeysFileName),
				Content: cfgKeysFile,
			})
	}

	p.GetFolder().Add(out...)

	return nil
}
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

type BuildProjectAction struct{}

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
