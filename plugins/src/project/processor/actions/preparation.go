package actions

import (
	"bytes"
	"github.com/Red-Sock/rscli/plugins/src/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/src/project/processor/patterns"
	"strings"

	"github.com/Red-Sock/rscli/pkg/folder"
)

func PrepareProjectStructure(p interfaces.Project) error {
	cmd := &folder.Folder{Name: "cmd"}

	cmd.Inner = append(cmd.Inner, &folder.Folder{
		Name: p.GetName(),
		Inner: []*folder.Folder{
			{
				Name:    "main.go",
				Content: patterns.MainFile,
			},
		},
	})
	fldr := p.GetFolder()
	fldr.Add(cmd)

	fldr.Add(&folder.Folder{Name: "config"})

	fldr.Add(&folder.Folder{Name: "internal"})

	fldr.Add(&folder.Folder{
		Name: "pkg",
		Inner: []*folder.Folder{
			{Name: "swagger"},
			{Name: "api"},
		},
	})

	return nil
}

func PrepareConfigFolders(p interfaces.Project) error {
	cfg := p.GetConfig()

	configFolders := make([]*folder.Folder, 0, 1)

	dsFolders, err := cfg.ExtractDataSources()
	if err != nil {
		return err
	}
	if dsFolders != nil {
		configFolders = append(configFolders, dsFolders)
	}

	p.GetFolder().AddWithPath([]string{"internal"}, configFolders...)
	return nil
}

func PrepareAPIFolders(p interfaces.Project) error {
	cfg := p.GetConfig()

	serverFolders, err := cfg.ExtractServerOptions()
	if err != nil {
		return err
	}

	if serverFolders == nil {
		return nil
	}

	projMainFile := p.GetFolder().GetByPath("cmd", p.GetName(), "main.go")

	importReplace := make([]string, 0, len(serverFolders.Inner))
	serversInit := make([]string, 0, len(serverFolders.Inner))

	for _, serv := range serverFolders.Inner {
		importReplace = append(importReplace, serv.Name+" \""+p.GetName()+"/internal/transport/"+serv.Name+"\"")
		serversInit = append(serversInit, "mngr.AddServer("+serv.Name+".NewServer(cfg))")
	}

	projMainFile.Content = bytes.ReplaceAll(
		projMainFile.Content,
		[]byte("//_transport_imports"),
		[]byte(strings.Join(importReplace, "\n\t")))

	projMainFile.Content = bytes.ReplaceAll(
		projMainFile.Content,
		[]byte("//_initiation_of_servers"),
		[]byte(strings.Join(serversInit, "\n\t")))

	serverFolders.Inner = append(serverFolders.Inner, &folder.Folder{
		Name:    "manager.go",
		Content: patterns.ServerManagerPattern,
	})

	if serverFolders != nil {
		p.GetFolder().AddWithPath([]string{"internal"}, serverFolders)
	}

	return nil
}

func PrepareExamplesFolders(p interfaces.Project) error {
	p.GetFolder().Add(&folder.Folder{
		Name: "examples",
		Inner: []*folder.Folder{
			{
				Name:    "api.http",
				Content: patterns.ApiHTTP,
			},
			{
				Name:    "http-client.env.json",
				Content: patterns.HttpEnvironment,
			},
		},
	})
	return nil
}

func PrepareEnvironmentFolders(p interfaces.Project) error {
	p.GetFolder().Add(
		[]*folder.Folder{
			{
				Name:    "Dockerfile",
				Content: patterns.Dockerfile,
			},
			{
				Name:    "README.md",
				Content: bytes.ReplaceAll(patterns.Readme, []byte("{{PROJECT_NAME}}"), []byte(p.GetName())),
			},
		}...,
	)
	return nil
}
