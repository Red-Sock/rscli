package go_actions

import (
	"path"
	"text/template"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

var configKeysTemplate *template.Template

func init() {
	var err error
	configKeysTemplate, err = template.New("configTemplate").Funcs(template.FuncMap{
		"SnakeToPascal": cases.SnakeToPascal,
	}).Parse(`
package config

const ({{range $_, $val := .}}
	{{ SnakeToPascal $val.Name }} = "{{ $val.Name }}"{{end}}
)
`)
	if err != nil {
		panic(errors.Wrap(err, "error parsing config keys template"))
	}
}

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p interfaces.Project) (err error) {
	configFolderPath := path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder)
	p.GetFolder().Add(&folder.Folder{
		Name: configFolderPath,
	})

	goConfigFolder := p.GetFolder().GetByPath(configFolderPath)

	err = a.generateGoKeysFile(p, goConfigFolder)
	if err != nil {
		return errors.Wrap(err, "error generating keys go-file")
	}

	a.generateGoConfigFiles(goConfigFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return errors.Wrap(err, "error generating config yaml-files")
	}

	return nil
}
func (a PrepareGoConfigFolderAction) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareGoConfigFolderAction) generateGoKeysFile(p interfaces.Project, goConfigFolder *folder.Folder) error {
	keys, err := matreshka.GenerateKeys(*p.GetConfig().AppConfig)
	if err != nil {
		return errors.Wrap(err, "error generating environment keys")
	}

	if len(keys) == 0 {
		goConfigFolder.GetByPath(projpatterns.ConfigKeysFileName).Delete()
		return nil
	}

	buf := &rw.RW{}
	err = configKeysTemplate.Execute(buf, keys)
	if err != nil {
		return errors.Wrap(err, "error executing config keys template on generating")
	}

	cfgKeysFile := &folder.Folder{
		Name:    projpatterns.ConfigKeysFileName,
		Content: buf.Bytes(),
	}

	goConfigFolder.Add(cfgKeysFile)

	return nil
}

func (a PrepareGoConfigFolderAction) generateGoConfigFiles(goConfigFolder *folder.Folder) {
	goConfigFolder.Add(projpatterns.ConfigFile.Copy())
	goConfigFolder.Add(projpatterns.AutoloadConfigFile.Copy())
	goConfigFolder.Add(projpatterns.StaticConfigFile.Copy())
}

func (a PrepareGoConfigFolderAction) generateConfigYamlFile(p interfaces.Project) error {
	b, err := p.GetConfig().Marshal()
	if err != nil {
		return err
	}

	configFolder := p.GetFolder().GetByPath(projpatterns.ConfigsFolder)

	configFolder.Add(&folder.Folder{
		Name:    projpatterns.ConfigTemplate,
		Content: b,
	})
	configFolder.Add(&folder.Folder{
		Name:    projpatterns.DevConfigYamlFile,
		Content: b,
	})

	return nil
}
