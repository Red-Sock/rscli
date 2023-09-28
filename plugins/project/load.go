package project

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
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

func LoadProject(pth string, cfg *rscliconfig.RsCliConfig) (*Project, error) {
	configDirPath := path.Join(pth, path.Dir(cfg.Env.PathToConfig))

	dir, err := os.ReadDir(configDirPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config folder")
	}

	var yamlFiles = map[string]struct{}{}

	for _, d := range dir {
		if strings.HasSuffix(d.Name(), ".yaml") {

			yamlFiles[d.Name()] = struct{}{}
		}
	}

	var configPath string
	for _, d := range configOrder {
		if _, ok := yamlFiles[d]; ok {
			configPath = d
			break
		}
	}

	configPath = path.Join(configDirPath, configPath)

	c, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config")
	}

	f, err := folder.Load(pth)
	if err != nil {
		return nil, err
	}

	modName := c.ExtractName()

	goModFile := f.GetByPath(patterns.GoMod)
	moduleBts := goModFile.Content[:bytes.IndexByte(goModFile.Content, '\n')]
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
		pthToCfg:    configPath,
		root:        f,
	}

	err = interfaces.LoadProjectVersion(p)
	if err != nil {
		return p, errors.Wrap(err, "error loading project version")
	}

	return p, nil
}