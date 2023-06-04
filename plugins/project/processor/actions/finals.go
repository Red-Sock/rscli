package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	config "github.com/Red-Sock/rscli/plugins/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/tidy"

	"github.com/Red-Sock/rscli/pkg/cmd"
	configpattern "github.com/Red-Sock/rscli/plugins/config/pkg/configstructs"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

const goBin = "go"

type InitGoModAction struct{}

func (a InitGoModAction) Do(p interfaces.Project) error {

	command := exec.Command(goBin, "mod", "init", p.GetProjectModName())

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

	return nil
}
func (a InitGoModAction) NameInAction() string {
	return "Initiating go project"
}

type MoveCfgAction struct{}

func (a MoveCfgAction) Do(p interfaces.Project) error {

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

	cfg.AppInfo.Name = p.GetProjectModName()
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
func (a MoveCfgAction) NameInAction() string {
	return "Moving config"
}

type FixupProjectAction struct{}

func (a FixupProjectAction) Do(p interfaces.Project) error {

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
func (a FixupProjectAction) NameInAction() string {
	return "Performing project fix up"
}

type TidyAction struct{}

func (a TidyAction) Do(p interfaces.Project) error {
	err := tidy.Api(p)
	if err != nil {
		return errors.Wrap(err, "error during api tiding")
	}

	err = tidy.Config(p)
	if err != nil {
		return errors.Wrap(err, "error during config tiding")
	}

	err = tidy.DataSources(p)
	if err != nil {
		return errors.Wrap(err, "error during data source tiding")
	}

	renamer.ReplaceProjectName(p.GetProjectModName(), p.GetFolder())

	err = BuildConfigGoFolderAction{}.Do(p)
	if err != nil {
		return errors.Wrap(err, "error building go config folder")
	}

	return errors.Wrap(p.GetFolder().Build(), "error building project")
}
func (a TidyAction) NameInAction() string {
	return "Cleaning up the project"
}

// helping functions
