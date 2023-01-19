package embeded

import (
	"bytes"
	"fmt"
	"github.com/Red-Sock/rscli/internal/utils/shared"
	"github.com/Red-Sock/rscli/pkg/commands"
	"github.com/Red-Sock/rscli/pkg/rw"
	"github.com/pkg/errors"
	"io"
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

	println("building cmd plugin...")

	err = p.buildPluginCmd(pathToRepo)
	if err != nil {
		return err
	}

	println("cmd plugin built!")

	println("building ui plugin...")
	err = p.buildPluginUI(pathToRepo)
	if err != nil {
		return err
	}
	println("ui plugin built!")

	_, _ = os.Stdout.WriteString("plugin is successfully installed")

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
	println("Cloning git repository repoURL:\n")

	pluginDir := path.Join(allPluginsDir, URL.Host, URL.Path)

	_, err = os.ReadDir(pluginDir)
	if err == nil {
		return "", fmt.Errorf("%s is already installed. Delete it with %s %s %s and try again", repoURL, commands.RsCLI(), commands.Delete, repoURL)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", errors.Wrapf(err, "error can't perfom ReadDir")
	}

	repoPluginDir := path.Join(pluginDir, gitRepoTempNameDir)
	err = os.MkdirAll(repoPluginDir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "error creating directory for plugin repo")
	}

	err = p.gitFetch(repoPluginDir, repoURL)
	if err != nil {
		return "", err
	}

	version, err := p.getVersion(repoPluginDir)
	if err != nil {
		return "", err
	}

	println("Cloned successfully. Current version is " + version + ". Executing go mod...\n")

	err = p.gomod(repoPluginDir)
	if err != nil {
		return "", err
	}

	println("go mod executed!\n")
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

func (p *GetPlugin) gitFetch(dirPath, repoURL string) error {
	cmd := exec.Command("git", "clone", repoURL, ".")
	cmd.Dir = dirPath
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *GetPlugin) gomod(repoPluginDir string) error {
	gomodPath := path.Join(repoPluginDir, "go.mod")

	gomod, err := os.ReadFile(gomodPath)
	if err != nil {
		return errors.Wrapf(err, "error opening file %s", gomodPath)
	}

	startIdx := bytes.Index(gomod, []byte("github.com/Red-Sock/rscli"))
	oldImport := gomod[startIdx:]
	endIdx := bytes.Index(oldImport, []byte("\n"))
	oldImport = oldImport[:endIdx]

	bytes.ReplaceAll(gomod, oldImport, []byte("github.com/Red-Sock/rscli latest"))
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = repoPluginDir
	cmd.Stderr = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *GetPlugin) getVersion(repoPluginDir string) (string, error) {
	r := &rw.RW{}
	hashCommitCmd := exec.Command("git", "rev-parse", "HEAD")
	hashCommitCmd.Dir = repoPluginDir
	hashCommitCmd.Stderr = r

	commitHash, err := hashCommitCmd.Output()
	if err != nil {
		return "", err
	}

	tagCmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	tagCmd.Dir = repoPluginDir
	tagCmd.Stderr = r

	tag, err := tagCmd.Output()
	if err != nil {
		msg, err := io.ReadAll(r)
		if err != nil {
			return "", err
		}
		if string(msg) == "fatal: No names found, cannot describe anything." {
			return string(commitHash), nil
		}

		return "", errors.New(string(msg))
	}

	tagHashCmd := exec.Command("git", "show-ref", "-s", string(tag))
	tagHashCmd.Dir = repoPluginDir
	tagHashCmd.Stderr = r

	tagHash, err := tagHashCmd.Output()
	if err != nil {
		return "", err
	}

	if string(tagHash) == string(commitHash) {
		return string(tag), nil
	}

	return string(commitHash), nil
}
