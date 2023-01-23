package scripts

import (
	_ "embed"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/config"
)

const (
	EnvDir             = "environment"
	envExampleFile     = ".env.example"
	composeExampleFile = "docker-compose.yaml.example"
	makefileFile       = "Makefile"
)

const (
	projNamePattern     = "proj_name"
	projNameCapsPattern = "PROJ_NAME_CAPS"
)

var ErrEnvironmentExists = errors.New("environment already exists")

//go:embed files/.env
var envFile []byte

//go:embed files/docker-compose.yaml
var composeFile []byte

//go:embed files/Makefile
var makefile []byte

func RunCreate() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	err = createEnvDir()
	if err != nil {
		return err
	}

	var projs []string

	projs, err = getProjectsFromSubDir("./", cfg)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = createProjDirs(wd, projs)
	if err != nil {
		return err
	}

	return nil
}

func createEnvDir() error {
	_, err := os.ReadDir(EnvDir)
	if err == nil {
		return ErrEnvironmentExists
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	err = os.Mkdir(EnvDir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(EnvDir, envExampleFile), envFile, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(EnvDir, composeExampleFile), composeFile, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(EnvDir, makefileFile), selectMakefile(), 0755)
	if err != nil {
		return err
	}

	return nil
}

func getProjectsFromSubDir(pth string, cfg *config.Config) ([]string, error) {
	dirs, err := os.ReadDir(pth)
	if err != nil {
		return nil, err
	}

	projs := make([]string, 0, len(dirs))

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() && name != EnvDir {

			_, err = os.ReadFile(path.Join(name, strings.ReplaceAll(cfg.Env.PathToMain, projNamePattern, name)))
			if err != nil {
				if err == os.ErrNotExist {
					continue
				}
				return nil, err
			}

			projs = append(projs, d.Name())
		}
	}

	return projs, nil
}

func createProjDirs(projectsPath string, projcs []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	wd = path.Join(wd, EnvDir)
	for _, name := range projcs {
		projDir := path.Join(path.Join(projectsPath, EnvDir), name)
		err = os.Mkdir(projDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
