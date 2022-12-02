package project

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	configpattern "github.com/Red-Sock/rscli/pkg/config"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"gopkg.in/yaml.v3"
)

func initGoMod(p *Project) error {
	pth, ok := os.LookupEnv("GOROOT")
	if !ok {
		return fmt.Errorf("no go installed!\nhttps://golangr.com/install/")
	}

	cmd := exec.Command(pth+"/bin/go", "mod", "init", p.Name)
	wd, _ := os.Getwd()
	cmd.Dir = path.Join(wd, p.Name)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func prepareProjectStructure(p *Project) error {
	cmd := &Folder{name: "cmd"}
	cmd.inner = append(cmd.inner, &Folder{
		name: p.Name,
		inner: []*Folder{
			{
				name:    "main.go",
				content: mainFile,
			},
		},
	})
	p.f.inner = append(p.f.inner, cmd)

	p.f.inner = append(p.f.inner, &Folder{name: "config"})

	p.f.inner = append(p.f.inner, &Folder{name: "internal"})

	p.f.inner = append(p.f.inner, &Folder{
		name: "pkg",
		inner: []*Folder{
			{name: "swagger"},
			{name: "api"},
		},
	})

	return nil
}

func prepareConfigFolders(p *Project) error {
	if p.Cfg == nil {
		return nil
	}

	err := p.Cfg.parseSelf()
	if err != nil {
		return err
	}

	configFolders := make([]*Folder, 0, len(p.Cfg.values))

	dsFolders, err := p.Cfg.extractDataSources()
	if err != nil {
		return err
	}
	if dsFolders != nil {
		configFolders = append(configFolders, dsFolders)
	}

	serverFolders, err := p.Cfg.extractServerOptions()
	if err != nil {
		return err
	}
	if dsFolders != nil {
		configFolders = append(configFolders, serverFolders)
	}

	p.f.AddWithPath("internal", configFolders...)
	return nil
}

func prepareEnvironmentFolders(p *Project) error {
	p.f.inner = append(p.f.inner,
		[]*Folder{
			{
				name:    "Dockerfile",
				content: dockerfile,
			},
			{
				name:    "README.md",
				content: readme,
			},
		}...,
	)
	return nil
}

func buildConfigGoFolder(p *Project) error {
	out := []*Folder{
		{
			name: "config.go",
			content: []byte(
				strings.ReplaceAll(configurator, "{{projectNAME_}}", strings.ToUpper(p.Name)),
			),
		},
	}

	keys := generateConfigKeys(p.Name)
	if len(keys) != 0 {
		out = append(out,
			&Folder{
				name:    "keys.go",
				content: keys,
			})
	}

	p.f.AddWithPath(strings.Join(
		[]string{
			"internal",
			"config",
		}, string(os.PathSeparator)),
		out...,
	)

	return nil
}

func buildProject(p *Project) error {
	err := p.f.Build("")
	if err != nil {
		return err
	}
	return nil
}

func moveCfg(p *Project) error {
	if p.Cfg == nil {
		return nil
	}

	var content []byte

	oldPath := p.Cfg.path

	p.Cfg.path = path.Join(path.Dir(p.Cfg.path), p.Name, "config", config.FileName)

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

	err = os.WriteFile(p.Cfg.path, content, 0755)
	if err != nil {
		return err
	}

	return os.RemoveAll(oldPath)
}

func fetchDependencies(p *Project) error {
	pth, ok := os.LookupEnv("GOROOT")
	if !ok {
		return fmt.Errorf("no go installed!\nhttps://golangr.com/install/")
	}

	cmd := exec.Command(pth+"/bin/go", "mod", "tidy")
	wd, _ := os.Getwd()
	cmd.Dir = path.Join(wd, p.Name)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
