package project

import (
	"errors"
	"fmt"
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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

	flagAppName      = "name"
	flagAppNameShort = "n"

	flagCfgPath      = "cfg"
	flagCfgPathShort = "c"
)

type Project struct {
	Name string

	cfgPath string
}

func NewProject(args []string) (*Project, error) {
	return createProject(args)
}

func (p *Project) Create() error {
	folders, err := p.readConfig()
	if err != nil {
		return err
	}

	println(folders)

	return nil
}

func (p *Project) ValidateName() error {
	if strings.Contains(p.Name, " ") {
		return errors.New("name contains space symbols")
	}
	return nil
}

func (p *Project) readConfig() ([]folder, error) {
	if p.cfgPath == "" {
		return nil, nil
	}

	conf := make(map[string]interface{})

	bytes, err := os.ReadFile(p.cfgPath)
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

	return nil, err
}

func extractDataSources(ds map[string]interface{}) (folder, error) {
	out := folder{
		name: "data",
	}

	for dsn := range ds {
		out.inner = append(out.inner, folder{
			name: dsn,
		})
	}

	return out, nil
}

func createProject(args []string) (*Project, error) {
	p := &Project{}

	flags, err := utils.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	p.Name, err = extractNameFromFlags(flags)
	if err != nil {
		return nil, err
	}

	if p.Name == "" {
		p.cfgPath, err = tryFindConfig(flags)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func extractNameFromFlags(flagsArgs map[string][]string) (string, error) {
	name, ok := flagsArgs[flagAppName]
	if !ok {
		name, ok = flagsArgs[flagAppNameShort]
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

func extractCfgPathFromFlags(flagsArgs map[string][]string) (string, error) {
	name, ok := flagsArgs[flagCfgPath]
	if !ok {
		name, ok = flagsArgs[flagCfgPathShort]
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

func tryFindConfig(args map[string][]string) (string, error) {
	pth, err := extractCfgPathFromFlags(args)
	if err != nil {
		return "", err
	}

	if pth != "" {
		return pth, nil
	}

	currentDir := filepath.Dir("./")
	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		return "", err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = config.DefaultDir
			break
		}
	}

	if pth == "" {
		return "", ErrNoConfigNoAppNameFlag
	}

	return pth, nil
}

type folder struct {
	name  string
	inner []folder
}
