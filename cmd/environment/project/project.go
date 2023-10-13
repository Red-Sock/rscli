package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/makefile"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/cmd/environment/project/ports"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	pconfig "github.com/Red-Sock/rscli/plugins/project/config"
)

var ErrNoConfig = errors.New("no config found")

type globalEnvConfig interface {
	GetByName(envName string) (value string)
}

type Env struct {
	envDirPath string
	projPath   string

	Compose         *compose.Compose
	Environment     *env.Container
	ComposePatterns compose.PatternManager
	Config          *pconfig.Config
	Makefile        *makefile.Makefile

	globalEnvFile  globalEnvConfig
	globalMakefile *makefile.Makefile
	rscliConfig    *config.RsCliConfig
}

func LoadProjectEnvironment(
	cfg *config.RsCliConfig,
	envResourcePattern globalEnvConfig,
	globalMakefile *makefile.Makefile,
	pathToProjectEnv string,
	pathToProject string,
) (p *Env, err error) {
	p = &Env{
		envDirPath: pathToProjectEnv,
		projPath:   pathToProject,

		globalEnvFile:  envResourcePattern,
		globalMakefile: globalMakefile,
		rscliConfig:    cfg,
	}

	err = p.fetchComposeFile()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching compose file")
	}

	err = p.fetchEnvFile()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching .env file")
	}

	err = p.fetchConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching config")
	}

	err = p.fetchMakeFile()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching makefile")
	}

	return p, nil
}

func (e *Env) Tidy(pm *ports.PortManager, serviceEnabled bool) error {
	projName := path.Base(e.envDirPath)

	e.tidyConfigFile()

	e.tidyEnvFile()

	err := e.tidyResources(pm, projName, serviceEnabled)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on resources")
	}

	err = e.tidyServerAPIs(projName, pm)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on server api")
	}

	{
		pathToProjectEnvFile := path.Join(e.envDirPath, patterns.EnvFile.Name)

		envBytes := e.Environment.MarshalEnv()
		if len(envBytes) != 0 {
			err = io.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectName(envBytes, projName))
			if err != nil {
				return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
			}
		}
	}

	{
		var composeFile []byte
		composeFile, err = e.Compose.Marshal()
		if err != nil {
			return errors.Wrap(err, "error marshalling composer file")
		}

		pathToDockerComposeFile := path.Join(e.envDirPath, patterns.DockerComposeFile.Name)
		err = io.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectName(composeFile, projName))
		if err != nil {
			return errors.Wrap(err, "error writing docker compose file file")
		}
	}

	{
		err = e.Config.BuildTo(e.Config.GetPath())
		if err != nil {
			return errors.Wrap(err, "error writing env config")
		}
	}

	{
		e.tidyMakeFile(projName)

		var mkFile []byte
		mkFile, err = e.Makefile.Marshal()
		if err != nil {
			return errors.Wrap(err, "error marshalling makefile")
		}

		err = io.OverrideFile(path.Join(e.envDirPath, patterns.Makefile.Name), mkFile)
		if err != nil {
			return errors.Wrap(err, "error writing makefile")
		}
	}

	return nil
}

func (e *Env) fetchComposeFile() error {
	projectEnvComposeFilePath := path.Join(e.envDirPath, patterns.DockerComposeFile.Name)
	composeFile, err := os.ReadFile(projectEnvComposeFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project env docker-compose file "+projectEnvComposeFilePath)
		}
	}

	if len(composeFile) == 0 {
		globalEnvComposeFilePath := path.Join(path.Dir(e.envDirPath), patterns.DockerComposeFile.Name)
		composeFile, err = os.ReadFile(globalEnvComposeFilePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrap(err, "error reading global docker-compose file "+globalEnvComposeFilePath)
			}
		}
	}

	if len(composeFile) == 0 {
		projName := path.Base(e.envDirPath)
		composeFile = renamer.ReplaceProjectName(patterns.DockerComposeFile.Content, projName)
	}

	e.Compose, err = compose.NewComposeAssembler(composeFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}

func (e *Env) fetchEnvFile() error {
	dotEnvFilePath := path.Join(e.envDirPath, patterns.EnvFile.Name)
	envFile, err := os.ReadFile(dotEnvFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading project .env file "+dotEnvFilePath)
		}
	}

	e.Environment, err = env.NewEnvContainer(envFile)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	return nil
}
