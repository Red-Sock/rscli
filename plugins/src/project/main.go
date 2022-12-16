package main

import (
	"os"
	"path"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/project/config-processor/config"

	"github.com/Red-Sock/rscli/pkg/flag"

	projectprocessor "github.com/Red-Sock/rscli/pkg/service/project/project-processor"
)

const (
	FlagAppName      = "name"
	FlagAppNameShort = "n"

	FlagCfgPath      = "cfg"
	FlagCfgPathShort = "c"

	FlagProjectPath      = "project-path"
	FlagProjectPathShort = "p"
)

var Plug plugin

type plugin struct{}

func (p *plugin) GetName() string {
	return "project"
}

func (p *plugin) Run(args []string) error {
	projArgs := projectprocessor.CreateArgs{}

	flags := flag.ParseArgs(args)

	var err error
	// Define project name
	projArgs.Name, err = flag.ExtractOneValueFromFlags(flags, FlagAppName, FlagAppNameShort)
	if err != nil {
		return err
	}

	// Define path to configuration file
	projArgs.CfgPath, err = flag.ExtractOneValueFromFlags(flags, FlagCfgPath, FlagCfgPathShort)
	if err != nil {
		return err
	}
	if projArgs.CfgPath == "" {
		projArgs.CfgPath, err = findConfigPath()
		if err != nil {
			return err
		}
	}

	projArgs.ProjectPath, err = flag.ExtractOneValueFromFlags(flags, FlagProjectPath, FlagProjectPathShort)
	if projArgs.ProjectPath == "" {
		projArgs.ProjectPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	proj, err := projectprocessor.New(projArgs)
	if err != nil {
		return err
	}
	err = proj.Validate()
	if err != nil {
		return err
	}

	err = proj.Build()
	if err != nil {
		return err
	}

	return nil
}

func findConfigPath() (pth string, err error) {
	currentDir := "./"

	var dirs []os.DirEntry
	dirs, err = os.ReadDir(currentDir)
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
		return "", nil
	}

	confs, err := os.ReadDir(pth)
	if err != nil {
		return "", err
	}
	for _, f := range confs {
		name := f.Name()
		if strings.HasSuffix(name, config.FileName) {
			pth = path.Join(pth, name)
			break
		}
	}

	return pth, nil
}
