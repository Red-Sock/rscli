package scripts

import (
	_ "embed"
	"io/fs"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/config"
)

const (
	EnvDir              = "environment"
	EnvFile             = ".env"
	envExampleFile      = ".env.example"
	composeExampleFile  = "docker-compose.example.yaml"
	dockerComposeFile   = "docker-compose.yaml"
	makefileExampleFile = "Makefile"
)

const (
	projNamePattern               = "proj_name"
	projNameCapsPattern           = "PROJ_NAME_CAPS"
	datasourceCapsPostgresPattern = "DS_POSTGRES_NAME_CAPS"
)

var ErrEnvironmentExists = errors.New("environment already exists")

//go:embed patterns/files/.env
var envFile []byte

//go:embed patterns/files/docker-compose.yaml
var composeFile []byte

//go:embed patterns/files/Makefile
var makefile []byte

func RunCreate() error {
	cfg, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		return err
	}

	err = createEnvDir()
	if err != nil {
		return err
	}

	var projects []string
	projects, err = ListProjects(wd, cfg)
	if err != nil {
		return err
	}

	err = CreateEnvFoldersForProjects(wd, projects)
	if err != nil {
		return err
	}

	return RunSetUp(projects)
}

func CreateEnvFoldersForProjects(projectsPath string, projects []string) error {
	for _, name := range projects {
		projDir := path.Join(path.Join(projectsPath, EnvDir), name)

		err := os.Mkdir(projDir, os.ModePerm)
		if err != nil {
			return err
		}
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

	{
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

		err = os.WriteFile(path.Join(EnvDir, makefileExampleFile), selectMakefile(), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
