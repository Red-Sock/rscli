package config_generators

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
)

type loadConfigFileGenArgs struct {
	Configs []InternalConfig
}

type InternalConfig struct {
	FieldName    string
	StructName   string
	From         string
	ErrorMessage string
}

type internalConfigGenerator func() (InternalConfig, *folder.Folder, error)

func GenerateConfigFolder(cfg *config.Config) (*folder.Folder, error) {
	args := loadConfigFileGenArgs{}

	configFolder := &folder.Folder{}

	generators := make([]internalConfigGenerator, 0, 3)

	// Data sources
	if len(cfg.DataSources) != 0 {
		generators = append(generators, newGenerateDataSourcesConfigStruct(cfg.DataSources))
	}

	// Environment
	if len(cfg.Environment) != 0 {
		generators = append(generators, newGenerateEnvironmentConfigStruct(cfg.Environment))
	}

	for _, g := range generators {
		ic, f, err := g()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		configFolder.Add(f)
		args.Configs = append(args.Configs, ic)
	}

	autoLoadFile := &rw.RW{}
	err := configAutoLoadTemplate.Execute(autoLoadFile, args)
	if err != nil {
		return nil, errors.Wrap(err, "error generating load-config file ")
	}

	configFolder.Add(&folder.Folder{
		Name:    projpatterns.ConfigLoadFileName,
		Content: autoLoadFile.Bytes(),
	})
	return configFolder, nil
}
