package processor

import (
	"github.com/Red-Sock/rscli/pkg/cmd"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var configOrder = []string{
	"dev.yaml",
	"stage.yaml",
	"prod.yaml",
	"*.yaml",
}

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

	c, err := config.NewProjectConfig(path.Join(pth, patterns.ConfigsFolder, configPath))
	if err != nil {
		return nil, err
	}

	err = c.ParseSelf()
	if err != nil {
		return nil, err
	}

	name, err := c.ExtractName()
	if err != nil {
		return nil, err
	}

	f, err := folder.Load(pth)
	if err != nil {
		return nil, err
	}

	version, err := GetProjectVersion(pth)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:         name,
		ProjectPath:  pth,
		Cfg:          c,
		F:            f,
		RscliVersion: version,
	}, nil
}

func GetProjectVersion(wd string) (string, error) {
	out, err := cmd.Execute(cmd.Request{
		Tool:    "make",
		Args:    []string{"-f", patterns.RsCliMkFileName, "rscli-version"},
		WorkDir: wd,
	})
	if err != nil {
		return "", errors.Wrap(err, "error executing make rscli-version")
	}

	return out, nil
}
