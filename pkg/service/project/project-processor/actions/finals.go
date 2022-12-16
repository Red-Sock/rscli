package actions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/Red-Sock/rscli/pkg/service/project/config-processor/config"

	"github.com/Red-Sock/rscli/pkg/service/project/project-processor/interfaces"

	"github.com/Red-Sock/rscli/pkg/folder"

	configpattern "github.com/Red-Sock/rscli/pkg/config"
	"gopkg.in/yaml.v3"
)

func InitGoMod(p interfaces.Project) error {
	pth, ok := os.LookupEnv("GOROOT")
	if !ok {
		return fmt.Errorf("no go installed!\nhttps://golangr.com/install/")
	}

	cmd := exec.Command(pth+"/bin/go", "mod", "init", p.GetName())
	wd, _ := os.Getwd()
	cmd.Dir = path.Join(wd, p.GetName())
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func MoveCfg(p interfaces.Project) error {

	var content []byte

	sourceConfPath := p.GetConfig().GetPath()

	if sourceConfPath == "" {
		return nil
	}

	projectConfPath := path.Join(p.GetProjectPath(), "config", config.FileName)

	if sourceConfPath == projectConfPath {
		return nil
	}

	content, err := os.ReadFile(sourceConfPath)
	if err != nil {
		return err
	}

	var cfg configpattern.Config
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config from file %w", err)
	}

	cfg.AppInfo.Name = p.GetName()
	cfg.AppInfo.Version = "0.0.1"

	content, err = yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config into file %w", err)
	}

	err = os.WriteFile(projectConfPath, content, 0755)
	if err != nil {
		return err
	}

	p.GetConfig().SetPath(projectConfPath)

	return os.RemoveAll(sourceConfPath)
}

func FixupProject(p interfaces.Project) error {
	pth, ok := os.LookupEnv("GOROOT")
	if !ok {
		return fmt.Errorf("no go installed!\nhttps://golangr.com/install/")
	}

	wd, _ := os.Getwd()
	wd = path.Join(wd, p.GetName())

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

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		f.Content = bytes.ReplaceAll(f.Content, []byte("financial-microservice"), []byte(name))
		return
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
	}
}
