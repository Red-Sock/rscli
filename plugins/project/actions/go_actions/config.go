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
	configKeysTemplate, err = template.New("configTemplate").
		Parse(`
package config

const ({{range $_, $val := .}}
	{{ $val.Name }} = "{{ $val.Key }}"{{end}}
)
`)
	if err != nil {
		panic(errors.Wrap(err, "error parsing config keys template"))
	}
}

type PrepareGoConfigFolderAction struct{}

func (a PrepareGoConfigFolderAction) Do(p interfaces.Project) (err error) {
	goConfigFolderPath := path.Join(projpatterns.InternalFolder, projpatterns.ConfigsFolder)
	p.GetFolder().Add(&folder.Folder{
		Name: goConfigFolderPath,
	})
	goConfigFolder := p.GetFolder().GetByPath(goConfigFolderPath)

	err = a.generateGoKeysFile(p, goConfigFolder)
	if err != nil {
		return errors.Wrap(err, "error generating keys go-file")
	}

	goConfigFolder.Add(projpatterns.AutoloadConfigFile.Copy())

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
	matreshkaKeys := matreshka.GenerateKeys(p.GetConfig().AppConfig)

	type ConfigKey struct {
		Name string
		Key  string
	}

	keys := make([]ConfigKey, 0, len(matreshkaKeys.DataSources)+len(matreshkaKeys.Servers)+len(matreshkaKeys.Environment))

	if len(keys) == 0 {
		goConfigFolder.GetByPath(projpatterns.ConfigKeysFileName).Delete()
		return nil
	}

	for _, k := range matreshkaKeys.DataSources {
		keys = append(keys, ConfigKey{
			Name: "Resource" + cases.SnakeToPascal(k),
		})
	}
	for _, k := range matreshkaKeys.Servers {
		keys = append(keys, ConfigKey{
			Name: "Server" + cases.SnakeToPascal(k),
		})
	}

	for _, k := range matreshkaKeys.Environment {
		keys = append(keys, ConfigKey{
			Name: cases.SnakeToPascal(k),
		})
	}

	buf := &rw.RW{}
	err := configKeysTemplate.Execute(buf, keys)
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
