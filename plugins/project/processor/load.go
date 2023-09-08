package processor

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

const (
	DevConfigFileName  = "dev.yaml"
	StgConfigFileName  = "stage.yaml"
	ProdConfigFileName = "prod.yaml"
)

var configOrder = []string{
	DevConfigFileName,
	StgConfigFileName,
	ProdConfigFileName,
}

// TODO
func LoadProject(pth string) (*Project, error) {
	dir, err := os.ReadDir(path.Join(pth, patterns.ConfigsFolder))
	if err != nil {
		return nil, errors.Wrap(err, "error reading config folder")
	}

	var yamls = map[string]struct{}{}
	var configPath string
	for _, d := range dir {
		if strings.HasSuffix(d.Name(), ".yaml") {
			configPath = d.Name()
			yamls[configPath] = struct{}{}
		}
	}

	for _, d := range configOrder {
		if _, ok := yamls[d]; ok {
			configPath = d
			break
		}
	}

	c, err := config.ReadConfig(path.Join(pth, patterns.ConfigsFolder, configPath))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config")
	}

	f, err := folder.Load(pth, "")
	if err != nil {
		return nil, err
	}

	modName := c.ExtractName()

	gomodF := f.GetByPath(patterns.GoMod)
	moduleBts := gomodF.Content[:bytes.IndexByte(gomodF.Content, '\n')]
	moduleBts = moduleBts[1+bytes.IndexByte(moduleBts, ' '):]

	if modName != string(moduleBts) {
		modName = string(moduleBts)
	}

	name := modName

	if nameStartIdx := strings.LastIndex(modName, "/"); nameStartIdx != -1 {
		name = modName[nameStartIdx+1:]
	}

	p := &Project{
		Name:        name,
		ProjectPath: pth,
		Cfg:         c,
		root:        f,
	}

	err = interfaces.LoadProjectVersion(p)
	if err != nil {
		return p, errors.Wrap(err, "error loading project version")
	}

	return p, nil
}
