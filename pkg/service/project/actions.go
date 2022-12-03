package project

import (
	"bytes"
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

	p.f.AddWithPath([]string{"internal"}, configFolders...)
	return nil
}

func prepareAPIFolders(p *Project) error {
	serverFolders, err := p.Cfg.extractServerOptions()
	if err != nil {
		return err
	}

	projMainFile := p.f.GetByPath("cmd", p.Name, "main.go")

	importReplace := make([]string, 0, len(serverFolders.inner))
	serversInit := make([]string, 0, len(serverFolders.inner))

	for _, serv := range serverFolders.inner {
		importReplace = append(importReplace, serv.name+" \""+p.Name+"/internal/transport/"+serv.name+"\"")
		serversInit = append(serversInit, "mngr.AddServer("+serv.name+".NewServer(cfg))")
	}

	projMainFile.content = bytes.ReplaceAll(
		projMainFile.content,
		[]byte("//_transport_imports"),
		[]byte(strings.Join(importReplace, "\n\t")))

	projMainFile.content = bytes.ReplaceAll(
		projMainFile.content,
		[]byte("//_initiation_of_servers"),
		[]byte(strings.Join(serversInit, "\n\t")))

	serverFolders.inner = append(serverFolders.inner, &Folder{
		name:    "manager.go",
		content: managerPattern,
	})

	if serverFolders != nil {
		p.f.AddWithPath([]string{"internal"}, serverFolders)
	}

	return nil
}

func prepareExamplesFolders(p *Project) error {
	p.f.inner = append(p.f.inner, &Folder{
		name: "examples",
		inner: []*Folder{
			{
				name:    "api.http",
				content: apiHTTP,
			},
			{
				name:    "http-client.env.json",
				content: httpEnvironment,
			},
		},
	})
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

	keys, err := generateConfigKeys(p.Name, p.Cfg.path)
	if err != nil {
		return err
	}
	if len(keys) != 0 {
		out = append(out,
			&Folder{
				name:    "keys.go",
				content: bytes.ReplaceAll(configKeys, []byte("// _config_keys_goes_here"), keys),
			})
	}

	p.f.AddWithPath(
		[]string{
			"internal",
			"config",
		},
		out...,
	)

	return nil
}

func buildProject(p *Project) error {

	changeProjectName(p.Name, &p.f)

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

func fixupProject(p *Project) error {
	pth, ok := os.LookupEnv("GOROOT")
	if !ok {
		return fmt.Errorf("no go installed!\nhttps://golangr.com/install/")
	}

	wd, _ := os.Getwd()
	wd = path.Join(wd, p.Name)

	cmd := exec.Command(pth+"/bin/go", "mod", "tidy")
	cmd.Dir = wd
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command(pth+"/bin/go", "fmt", "./...")
	cmd.Dir = wd
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// helping functions

func changeProjectName(name string, f *Folder) {
	if f.content != nil {
		f.content = bytes.ReplaceAll(f.content, []byte("financial-microservice"), []byte(name))
		return
	}
	for _, innerFolder := range f.inner {
		changeProjectName(name, innerFolder)
	}
}
