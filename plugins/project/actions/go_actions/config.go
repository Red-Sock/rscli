package go_actions

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.verv.tech/matreshka"
	"go.verv.tech/matreshka/environment"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/config_generators"
)

const (
	LogLevelEvonName = "log-level"
	LogLevelTrace    = "Trace"
	LogLevelDebug    = "Debug"
	LogLevelInfo     = "Info"
	LogLevelWarn     = "Warn"
	LogLevelError    = "Error"
	LogLevelFatal    = "Fatal"
	LogLevelPanic    = "Panic"

	LogFormatEvonName = "log-format"
	LogFormatJSON     = "JSON"
	LogFormatTEXT     = "TEXT"
)

type GenerateProjectConfig struct {
}

func (a GenerateProjectConfig) Do(p project.IProject) error {
	cfg := p.GetConfig()

	envVars := map[string]*environment.Variable{}

	for _, v := range cfg.Environment {
		envVars[v.Name] = v
	}

	if envVars[LogLevelEvonName] == nil {
		cfg.Environment = append(cfg.Environment, &environment.Variable{
			Name: LogLevelEvonName,
			Type: environment.VariableTypeStr,
			Enum: []any{
				LogLevelTrace,
				LogLevelDebug,
				LogLevelInfo,
				LogLevelWarn,
				LogLevelError,
				LogLevelFatal,
				LogLevelPanic,
			},
			Value: LogLevelInfo,
		})
	}

	if envVars[LogFormatEvonName] == nil {
		cfg.Environment = append(cfg.Environment, &environment.Variable{
			Name: LogFormatEvonName,
			Type: environment.VariableTypeStr,
			Enum: []any{
				LogFormatJSON,
				LogFormatTEXT,
			},
			Value: LogFormatTEXT,
		})
	}

	return nil
}
func (a GenerateProjectConfig) NameInAction() string {
	return "Generating project config"
}

type PrepareConfigFolder struct{}

func (a PrepareConfigFolder) Do(p project.IProject) (err error) {
	cfgFolder, err := config_generators.GenerateConfigFolder(p.GetConfig())
	if err != nil {
		return rerrors.Wrap(err, "error generating config folder")
	}

	cfgFolder.Name = path.Join(patterns.InternalFolder, patterns.ConfigsFolder)

	p.GetFolder().Add(cfgFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return rerrors.Wrap(err, "error generating config yaml-files")
	}

	return nil
}
func (a PrepareConfigFolder) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareConfigFolder) generateConfigYamlFile(p project.IProject) error {
	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)

	newConfig := p.GetConfig()

	for _, cfgName := range []string{
		patterns.ConfigDevYamlFile,
		patterns.ConfigTemplateYaml,
		patterns.ConfigMasterYamlFile,
	} {
		err := appendToConfig(newConfig.AppConfig, configFolder, cfgName)
		if err != nil {
			return rerrors.Wrap(err, "error appending changes to dev config")
		}
	}

	return nil
}

func appendToConfig(newConfig matreshka.AppConfig, configFolder *folder.Folder, path string) (err error) {
	currentConfig := matreshka.NewEmptyConfig()

	configFile := configFolder.GetByPath(path)
	if configFile == nil {
		configFile = &folder.Folder{
			Name: path,
		}
		configFolder.Add(configFile)
	}

	if len(configFile.Content) != 0 {
		err = currentConfig.Unmarshal(configFile.Content)
		if err != nil {
			return rerrors.Wrap(err, "error reading dev config file")
		}
	}

	currentConfig = matreshka.MergeConfigs(currentConfig, newConfig)

	configFile.Content, err = currentConfig.Marshal()
	if err != nil {
		return rerrors.Wrap(err, "error marshalling dev config to yaml")
	}

	return nil
}
