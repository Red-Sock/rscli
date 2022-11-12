package project

import (
	"errors"
	"fmt"
	"github.com/Red-Sock/rscli/internal/utils"
	configpattern "github.com/Red-Sock/rscli/pkg/config"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

var commands = []string{"p", "project"}

func Command() []string {
	return commands
}

var (
	ErrNoConfigNoAppNameFlag = errors.New("no config or app name flag was specified")

	ErrNoArgumentsSpecifiedForFlag = errors.New("flag specified but no name was given")
	ErrFlagHasTooManyArguments     = errors.New("too many arguments specified for flag")
)

const (
	flagCreate = "create"

	FlagAppName      = "name"
	FlagAppNameShort = "n"

	FlagCfgPath      = "cfg"
	FlagCfgPathShort = "c"
)

type Project struct {
	Name string

	CfgPath string
}

func NewProject(args []string) (Project, error) {
	return createProject(args)
}

func (p *Project) Create() error {
	folders, err := p.readConfig()
	if err != nil {
		return err
	}

	patterns := p.readPatterns()

	folders = append(folders, patterns...)

	folders = append(folders, folder{name: "config"})

	projFolder := folder{
		name:  "./" + p.Name,
		inner: folders,
	}

	err = projFolder.MakeAll("")
	if err != nil {
		return err
	}

	return p.moveCfg()

}

func (p *Project) ValidateName() error {
	if p.Name == "" {
		return errors.New("no name entered")
	}

	if strings.Contains(p.Name, " ") {
		return errors.New("name contains space symbols")
	}

	return nil
}

func (p *Project) readConfig() ([]folder, error) {
	if p.CfgPath == "" {
		return nil, nil
	}

	conf := make(map[string]interface{})

	bytes, err := os.ReadFile(p.CfgPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	dataSources, ok := conf[config.DataSourceKey]
	if !ok {
		return nil, nil
	}
	ds, ok := dataSources.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	out := make([]folder, 0, len(conf))

	dsFolders, err := extractDataSources(ds)
	if err != nil {
		return nil, err
	}

	out = append(out, dsFolders)

	return out, err
}

func (p *Project) moveCfg() error {
	if p.CfgPath == "" {
		return nil
	}

	var content []byte

	oldPath := p.CfgPath
	p.CfgPath = path.Join(p.Name, "config", config.FileName)

	content, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}

	var cfg configpattern.Config
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config from file %w", err)
	}

	cfg.AppInfo.Name = p.Name
	cfg.AppInfo.Version = "0.0.1"

	content, err = yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config into file %w", err)
	}

	err = os.WriteFile(p.CfgPath, content, 0755)
	if err != nil {
		return nil
	}

	return os.RemoveAll(oldPath)
}

func (p *Project) readPatterns() (out []folder) {
	cmd := folder{
		name:  "cmd",
		inner: nil,
	}

	cmd.inner = append(cmd.inner, folder{
		name: p.Name,
		inner: []folder{
			{
				name:    "main.go",
				content: mainFile,
			},
		},
	})
	out = append(out, cmd)

	out = append(out, folder{
		name: "internal",
	})

	out = append(out, folder{
		name: "pkg",
		inner: []folder{
			{
				name: "swagger",
			},
			{
				name: "api",
			},
		},
	})

	return out
}

func createProject(args []string) (Project, error) {
	p := Project{}

	flags, err := utils.ParseArgs(args)
	if err != nil {
		return p, err
	}

	p.Name, err = extractNameFromFlags(flags)
	if err != nil {
		return p, err
	}

	if p.Name == "" {
		err = p.tryFindConfig(flags)
		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func extractNameFromFlags(flagsArgs map[string][]string) (string, error) {
	name, ok := flagsArgs[FlagAppName]
	if !ok {
		name, ok = flagsArgs[FlagAppNameShort]
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
