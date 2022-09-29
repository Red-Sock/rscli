package project

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/Red-Sock/rscli/internal/service/config"
	"github.com/Red-Sock/rscli/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

//go:embed main.go.pattern
var mainFile []byte

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

	patterns, err := p.readPatterns()
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

	return out, err
}

func (p *Project) moveCfg() error {
	if p.cfgPath == "" {
		return nil
	}
	cfgDir, _ := path.Split(p.cfgPath)
	cfgs, err := os.ReadDir(cfgDir)
	if err != nil {
		return err
	}

	var content []byte

	for _, c := range cfgs {
		oldPath := path.Join(cfgDir, c.Name())
		content, err = os.ReadFile(oldPath)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(p.Name, "config", c.Name()), content, 0755)
		if err != nil {
			return nil
		}

		err = os.Remove(oldPath)
		if err != nil {
			return err
		}
	}

	return os.Remove(cfgDir)
}

func (p *Project) readPatterns() (out []folder, err error) {

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

	return out, nil
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

	currentDir := "./"
	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		return "", err
	}

	for _, d := range dirs {
		if d.Name() == config.DefaultDir {
			pth = path.Join(currentDir, config.DefaultDir)
			break
		}
	}

	if pth == "" {
		return "", ErrNoConfigNoAppNameFlag
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return "", err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, "config.yaml") {
			pth = path.Join(pth, name)
			break
		}
	}

	return pth, nil
}

type folder struct {
	name    string
	inner   []folder
	content []byte
}

func (f *folder) MakeAll(root string) error {
	pth := path.Join(root, f.name)

	if len(f.content) != 0 {
		fw, err := os.Create(pth)
		if err != nil {
			return err
		}
		defer fw.Close()
		_, err = fw.Write(f.content)
		return err
	}

	err := os.MkdirAll(pth, 0755)
	if err != nil {
		return err
	}

	for _, d := range f.inner {
		err = d.MakeAll(pth)
		if err != nil {
			return err
		}
	}
	return nil
}
