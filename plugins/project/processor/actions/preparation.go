package actions

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
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

	fldr.AddWithPath([]string{patterns.InternalFolder, patterns.UtilsFolder, patterns.CloserFolder}, &folder.Folder{
		Name:    patterns.CloserFile,
		Content: patterns.UtilsCloser,
	})

	return nil
}

func PrepareConfigFolders(p interfaces.Project) error {
	cfg := p.GetConfig()

	configFolders := make([]*folder.Folder, 0, 1)

	dsFolders, err := cfg.GetDataSourceFolders()
	if err != nil {
		return err
	}
	if dsFolders != nil {
		configFolders = append(configFolders, dsFolders)
	}

	p.GetFolder().AddWithPath([]string{"internal"}, configFolders...)
	return nil
}

func PrepareExamplesFolders(p interfaces.Project) error {

	if p.GetFolder().GetByPath("examples", "http-client.env.json") != nil {
		return nil
	}

	type envs struct {
		Dev       map[string]string `json:"dev"`
		DevDocker map[string]string `json:"dev-docker"`
	}
	var e = envs{
		Dev:       map[string]string{},
		DevDocker: map[string]string{},
	}

	servers, err := p.GetConfig().GetServerOptions()
	if err != nil {
		return err
	}

	for _, item := range servers {
		e.Dev[item.Name] = "0.0.0.0:" + strconv.FormatUint(uint64(item.Port), 10)
		e.DevDocker[item.Name] = "0.0.0.0:1" + strconv.FormatUint(uint64(item.Port), 10)
	}

	eB, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		return err
	}

	p.GetFolder().Add(&folder.Folder{
		Name: "examples",
		Inner: []*folder.Folder{
			{
				Name:    "api.http",
				Content: patterns.ApiHTTP,
			},
			{
				Name:    "http-client.env.json",
				Content: eB,
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
				Content: bytes.ReplaceAll(patterns.Readme, []byte("{{PROJECT_NAME}}"), []byte(p.GetProjectModName())),
			},
			{
				Name:    ".gitignore",
				Content: patterns.GitIgnore,
			},
			{
				Name:    ".golangci.yaml",
				Content: patterns.Linter,
			},
			{
				Name:    "rscli.mk",
				Content: patterns.RscliMK,
			},
		}...,
	)
	return nil
}
