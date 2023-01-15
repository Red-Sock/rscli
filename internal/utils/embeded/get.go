package embeded

import (
	"fmt"
	"github.com/Red-Sock/rscli/internal/utils/shared"
	"github.com/Red-Sock/rscli/pkg/commands"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"os/exec"
	"path"
)

const gitRepoTempNameDir = "temp"

type GetPlugin struct {
}

func (p *GetPlugin) Run(flgs map[string][]string) error {
	allPluginsDir := shared.GetPluginsDir(flgs)

	pathToRepo, err := p.clone(allPluginsDir, flgs)
	if err != nil {
		return err
	}

	defer p.clean(pathToRepo)

	err = p.buildPluginCmd(pathToRepo)
	if err != nil {
		return err
	}

	err = p.buildPluginUI(pathToRepo)
	if err != nil {
		return err
	}

	return nil
}

func (p *GetPlugin) GetName() string {
	return commands.GetUtil
}

func (p *GetPlugin) clone(allPluginsDir string, flgs map[string][]string) (string, error) {
	if len(flgs) != 1 {
		return "", fmt.Errorf("invalid amount of agruments for %s plugins. Expected %d got %d", commands.GetUtil, 1, len(flgs))
	}

	var repoURL string
	for k := range flgs {
		repoURL = k
	}

	URL, err := url.Parse(repoURL)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing url %s", repoURL)
	}

	pluginDir := path.Join(allPluginsDir, URL.Host, URL.Path)

	_, err = os.ReadDir(pluginDir)
	if err == nil {
		return "", fmt.Errorf("%s is already installed. Delete it with %s %s and try again", repoURL, commands.RsCLI(), commands.Delete)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", errors.Wrapf(err, "error can't perfom ReadDir")
	}
	repoPluginDir := path.Join(pluginDir, gitRepoTempNameDir)
	_, err = git.PlainClone(repoPluginDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout, // todo replace with framework stdout
	})
	if err != nil {
		return "", errors.Wrapf(err, "error cloning repository %s", repoURL)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = repoPluginDir
	cmd.Stderr = os.Stdout // todo replace with framework stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return pluginDir, nil
}

func (p *GetPlugin) buildPluginCmd(newPluginDir string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", path.Join(newPluginDir, "cmd.so"), "main.go")

	cmd.Dir = path.Join(newPluginDir, gitRepoTempNameDir)
	cmd.Stderr = os.Stdout // todo replace with framework stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *GetPlugin) buildPluginUI(newPluginDir string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", path.Join(newPluginDir, "ui.so"), "main.go")

	cmd.Dir = path.Join(newPluginDir, gitRepoTempNameDir, "ui")
	cmd.Stderr = os.Stdout // todo replace with framework stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *GetPlugin) clean(dirPath string) {
	err := os.RemoveAll(path.Join(dirPath, gitRepoTempNameDir))
	if err != nil {
		fmt.Printf("error cleaning up %s: %s\n", dirPath, err)
	}
}
