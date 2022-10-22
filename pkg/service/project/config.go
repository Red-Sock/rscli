package project

import (
	"fmt"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"os"
	"path"
	"strings"
)

func (p *Project) tryFindConfig(args map[string][]string) error {
	pth, err := extractCfgPathFromFlags(args)
	if err != nil {
		return err
	}

	if pth != "" {
		p.CfgPath = pth
		return nil
	}

	currentDir := "./"
	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		return err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = path.Join(currentDir, config.DefaultDir)
			break
		}
	}

	if pth == "" {
		return ErrNoConfigNoAppNameFlag
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, config.FileName) {
			pth = path.Join(pth, name)
			break
		}
	}
	p.CfgPath = pth
	return nil
}

func extractDataSources(ds map[string]interface{}) (folder, error) {
	out := folder{
		name: "data",
	}

	for dsn := range ds {
		out.inner = append(out.inner, folder{
			name: dsn,
		})
	}

	return out, nil
}

func extractCfgPathFromFlags(flagsArgs map[string][]string) (string, error) {
	name, ok := flagsArgs[FlagCfgPath]
	if !ok {
		name, ok = flagsArgs[FlagCfgPathShort]
		if !ok {
			return "", nil
		}
	}
	if len(name) == 0 {
		return "", fmt.Errorf("%w expected 1 got 0 ", ErrNoArgumentsSpecifiedForFlag)
	}

	if len(name) > 1 {
		return "", fmt.Errorf("%w expected 1 got %d", ErrFlagHasTooManyArguments, len(name))
	}

	return name[0], nil
}
