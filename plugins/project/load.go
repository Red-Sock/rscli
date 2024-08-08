package project

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

const (
	ProdConfigFileName = "config.yaml"
	DevConfigFileName  = "dev.yaml"
	StgConfigFileName  = "stage.yaml"
)

var configOrder = []string{
	DevConfigFileName,
	StgConfigFileName,
	ProdConfigFileName,
}

func LoadProject(pth string, cfg *rscliconfig.RsCliConfig) (*Project, error) {
	c, err := LoadProjectConfig(pth, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error loading project config")
	}

	root, err := folder.Load(pth)
	if err != nil {
		return nil, err
	}

	p := &Project{
		ProjectPath: pth,
		Cfg:         c,
		root:        *root,
	}

	projectLoaders := []func(p *Project) (name *string){
		goProjectLoader,
		unknownProjectLoader,
	}

	for _, pLoader := range projectLoaders {
		name := pLoader(p)
		if name != nil {
			p.Name = *name
		}
	}

	err = proj_interfaces.LoadProjectVersion(p)
	if err != nil {
		return p, errors.Wrap(err, "error loading project version")
	}

	return p, nil
}

func LoadProjectConfig(projectPath string, cfg *rscliconfig.RsCliConfig) (c *config.Config, err error) {
	c = &config.Config{}

	configDirPath := path.Join(projectPath, path.Dir(cfg.Env.PathToConfig))

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

	c.Path = path.Join(configDirPath, configPath)

	c.AppConfig, err = matreshka.ReadConfigs(c.Path)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config")
	}

	return c, nil
}

func goProjectLoader(p *Project) (name *string) {
	goModFile := p.root.GetByPath(projpatterns.GoMod)
	if goModFile == nil {
		return nil
	}
	moduleBts := goModFile.Content[:bytes.IndexByte(goModFile.Content, '\n')]
	moduleBts = moduleBts[1+bytes.IndexByte(moduleBts, ' '):]

	modName := string(moduleBts)

	p.projType = proj_interfaces.ProjectTypeGo

	return &modName
}
func unknownProjectLoader(p *Project) *string {
	name := p.Cfg.AppInfo.Name
	return &name
}
