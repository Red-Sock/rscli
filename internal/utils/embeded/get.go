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
	"plugin"
	"strings"
)

const (
	gitRepoTempNameDir = "temp"
)

var (
	errPackageVersion = errors.New("more fresh package version required")
)

type GetPlugin struct{}

func (p *GetPlugin) Run(flgs map[string][]string) error {
	allPluginsDir := shared.GetPluginsDir(flgs)

	pathToRepo, err := p.clone(allPluginsDir, flgs)
	if err != nil {
		return err
	}

	defer p.clean(pathToRepo)

	err = p.buildPlugin(pathToRepo)
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

	println("Cloned successfully. Current version is " + version)
	println("Executing go mod...")

	err = p.gomod(repoPluginDir)
	if err != nil {
		return "", err
	}
	println("go mod executed!\n")

	return pluginDir, nil
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

	r := &rw.RW{}
	cmd.Stderr = r

	err := cmd.Run()
	if err != nil {
		msg, rErr := io.ReadAll(r.GetReader())
		if err != nil {
			return rErr
		}
		return errors.Wrap(err, string(msg))
	}
	return nil
}

func (p *GetPlugin) gomod(repoPluginDir string) error {
	cmd := exec.Command("go", "mod", "tidy")

	cmd.Dir = repoPluginDir

	r := &rw.RW{}
	cmd.Stderr = r

	err := cmd.Run()
	if err != nil {
		msg, rErr := io.ReadAll(r.GetReader())
		if err != nil {
			return rErr
		}
		return errors.Wrap(err, string(msg))
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
		msg, rErr := io.ReadAll(r.GetReader())
		if err != nil {
			return "", rErr
		}

		if bytes.Contains(msg, []byte("fatal: No names found, cannot describe anything.")) {
			println("i am here!!!!!")
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

func (p *GetPlugin) buildPlugin(pluginDir string) (err error) {
	pathToGitDir := path.Join(pluginDir, gitRepoTempNameDir)
	// building main CMD plugin
	err = p.build(path.Join(pathToGitDir, "cmd"), pluginDir, "cmd.so")
	if err != nil {
		return err
	}

	// building main UI plugin
	err = p.build(path.Join(pathToGitDir, "ui"), pluginDir, "ui.so")
	if err != nil {
		return err
	}

	return nil
}

func (p *GetPlugin) build(gitDir, pluginDir, name string) error {
	pathToSo := path.Join(pluginDir, name)

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pathToSo, "main.go")
	cmd.Dir = gitDir

	r := &rw.RW{}
	cmd.Stderr = r

	err := cmd.Run()
	if err != nil {
		msg, rErr := io.ReadAll(r.GetReader())
		if rErr != nil {
			return rErr
		}
		return errors.Wrap(err, string(msg))
	}

	println("opening " + pathToSo)
	_, err = plugin.Open(pathToSo)
	if err == nil {
		return nil
	}

	errStr := err.Error()
	println(errStr)
	os.Exit(1)
	const diffVerError = "plugin was built with a different version of package"

	if !strings.Contains(errStr, diffVerError) {
		return err
	}

	packageNameIdx := strings.LastIndex(errStr, diffVerError) + len(diffVerError)

	packageName := errStr[packageNameIdx+1:]
	println("package name is " + packageName)
	switch {
	case strings.HasPrefix(packageName, uikitName):
		err = p.updateDependencies(gitDir, uikitName)
	case strings.HasPrefix(packageName, rscliName):
		err = p.updateDependencies(gitDir, rscliName)
	default:
		return errors.New("cannot compile plugin: " + errStr)
	}

	if err != nil {
		return errors.Wrap(err, "error updating dependency")
	}

	println("removing unsuccessful plugin compilation " + pathToSo)

	err = os.Remove(pathToSo)
	if err != nil {
		return errors.Wrap(err, "error removing unsuccessfully compiled plugin")
	}

	return p.build(gitDir, pluginDir, name)
}

const (
	uikitName = "github.com/Red-Sock/rscli-uikit"
	rscliName = "github.com/Red-Sock/rscli"
)

func (p *GetPlugin) updateDependencies(repoPluginDir, depName string) error {
	println("Updating dependency " + depName)

	repoPluginDir = path.Dir(repoPluginDir)

	gomodPath := path.Join(repoPluginDir, "go.mod")

	println("gomodPath is: " + gomodPath)

	gomod, err := os.ReadFile(gomodPath)
	if err != nil {
		return errors.Wrapf(err, "error opening file %s", gomodPath)
	}
	gomod = findAndReplaceToLatest(gomod, depName+" ")

	err = os.WriteFile(gomodPath, gomod, 0755)
	if err != nil {
		return err
	}

	return p.gomod(repoPluginDir)
}

func findAndReplaceToLatest(src []byte, name string) []byte {
	startIdx := bytes.Index(src, []byte(name))
	oldImport := src[startIdx:]
	endIdx := bytes.Index(oldImport, []byte("\n"))
	oldImport = oldImport[:endIdx]
	src = bytes.ReplaceAll(src, oldImport, []byte(name+" latest"))

	return src
}
