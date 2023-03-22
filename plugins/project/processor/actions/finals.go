package actions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/plugins/project/processor/actions/tidy"

	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/pkg/folder"
	configpattern "github.com/Red-Sock/rscli/plugins/config/pkg/structs"
	config "github.com/Red-Sock/rscli/plugins/config/processor"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

const goBin = "go"

func InitGoMod(p interfaces.Project) error {

	command := exec.Command(goBin, "mod", "init", p.GetName())

	command.Dir = p.GetProjectPath()
	err := command.Run()
	if err != nil {
		return err
	}

	goMod, err := os.OpenFile(path.Join(p.GetProjectPath(), "go.mod"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer func() {
		err2 := goMod.Close()
		if err2 != nil {
			if err != nil {
				err = errors.Wrap(err, "error on closing"+err2.Error())
			} else {
				err = err2
			}
		}
	}()

	_, err = goMod.Write([]byte("\n// built via rscli v0.0.0"))
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

	wd, _ := os.Getwd()
	wd = path.Join(wd, p.GetName())

	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: wd,
	})
	if err != nil {
		return err
	}

	_, err = cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"fmt", "./..."},
		WorkDir: wd,
	})
	if err != nil {
		return err
	}

	return nil
}

func Tidy(p interfaces.Project) error {
	err := tidy.Api(p)
	if err != nil {
		return err
	}

	ReplaceProjectName(p.GetName(), p.GetFolder())

	err = BuildConfigGoFolder(p)
	if err != nil {
		return err
	}

	return p.GetFolder().Build()
}

// helping functions

const ProjectNamePattern = "financial-microservice"

func ReplaceProjectName(name string, f *folder.Folder) {
	if f.Content != nil {
		if idx := bytes.Index(f.Content, []byte(ProjectNamePattern)); idx != -1 {
			f.Content = bytes.ReplaceAll(f.Content, []byte(ProjectNamePattern), []byte(name))
			return
		}
	}
	for _, innerFolder := range f.Inner {
		ReplaceProjectName(name, innerFolder)
		if f.Name == ProjectNamePattern && len(f.Inner) == 0 {
			f.Name = name
		}
	}
}
