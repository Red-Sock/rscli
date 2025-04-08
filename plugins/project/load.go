package project

import (
	"bytes"
	"os"
	"path"
	"sort"
	"strings"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

const (
	prodConfigFileName     = "config.yaml"
	templateConfigFileName = "config_template.yaml"
	devConfigFileName      = "dev.yaml"
)

var configOrder = map[string]int{
	prodConfigFileName:     1,
	templateConfigFileName: 2,
	devConfigFileName:      3,
}

func LoadProject(pth string, cfg *rscliconfig.RsCliConfig) (*Project, error) {
	root, err := folder.Load(pth)
	if err != nil {
		return nil, err
	}

	conf, err := LoadProjectConfig(pth, cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "error loading project config")
	}

	p := &Project{
		Path: pth,
		Cfg:  conf,
		Root: *root,
	}

	projectLoaders := []func(p *Project) (name *string){
		unknownProjectLoader,
		goProjectLoader,
	}

	for _, pLoader := range projectLoaders {
		name := pLoader(p)
		if name != nil {
			p.Name = *name
		}
	}

	return p, nil
}

func LoadProjectConfig(projectPath string, cfg *rscliconfig.RsCliConfig) (c *config.Config, err error) {
	c = &config.Config{}

	c.ConfigDir = path.Join(projectPath, path.Dir(cfg.Env.PathToConfig))

	dir, err := os.ReadDir(c.ConfigDir)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading config folder")
	}

	configsPaths := make([]string, 0, 3)

	for _, d := range dir {
		if strings.HasSuffix(d.Name(), ".yaml") {
			_, ok := configOrder[d.Name()]
			if ok {
				configsPaths = append(configsPaths, path.Join(c.ConfigDir, d.Name()))
			}
		}
	}

	sort.Slice(configsPaths, func(i, j int) bool {
		return configOrder[configsPaths[i]] > configOrder[configsPaths[j]]
	})

	c.AppConfig, err = matreshka.ReadConfigs(configsPaths...) // TODO
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing config")
	}

	return c, nil
}

func goProjectLoader(p *Project) (name *string) {
	goModFile := p.Root.GetByPath(patterns.GoMod)
	if goModFile == nil {
		return nil
	}

	moduleBts := goModFile.Content[:bytes.IndexByte(goModFile.Content, '\n')]
	moduleBts = moduleBts[1+bytes.IndexByte(moduleBts, ' '):]

	modName := string(moduleBts)

	p.ProjType = TypeGo

	return &modName
}

func unknownProjectLoader(p *Project) *string {
	name := p.Cfg.AppInfo.Name
	return &name
}
