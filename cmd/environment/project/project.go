package project

import (
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
)

var ErrNoConfig = errors.New("no config found")

type globalEnvConfig interface {
	GetByName(envName string) (value string)
}

type ProjEnv struct {
	envProjPath string
	srcProjPath string

	Compose     envCompose
	Environment envVariables
	Makefile    envMakefile
	Config      envConfig

	ComposePatterns compose.PatternManager

	globalEnvFile  envVariables
	globalMakefile *makefile.Makefile
	rscliConfig    *config.RsCliConfig
}

func LoadProjectEnvironment(
	cfg *config.RsCliConfig,
	globalEnv *env.Container,
	globalMakefile *makefile.Makefile,

	pathToProjectEnv string,
	pathToProject string,

) (p *ProjEnv, err error) {
	p = &ProjEnv{
		envProjPath: pathToProjectEnv,
		srcProjPath: pathToProject,

		globalMakefile: globalMakefile,
		rscliConfig:    cfg,
	}

	err = p.Compose.fetch(pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching compose file")
	}

	err = p.Environment.fetch(globalEnv, pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching .env file")
	}

	err = p.Config.fetch(cfg, pathToProjectEnv, pathToProject)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching config")
	}

	err = p.Makefile.fetch(pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching makefile")
	}

	err = p.globalEnvFile.fetch(globalEnv, pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching makefile")
	}

	return p, nil
}

func (e *ProjEnv) Tidy(pm *ports.PortManager, serviceEnabled bool) error {
	projName := path.Base(e.envProjPath)

	err := e.tidyResources(pm, projName, serviceEnabled)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on resources")
	}

	err = e.tidyServerAPIs(projName, pm)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on server api")
	}

	err = e.flush(projName)

	return nil
}

func (e *ProjEnv) flush(projName string) (err error) {
	{
		pathToProjectEnvFile := path.Join(e.envProjPath, patterns.EnvFile.Name)

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

		pathToDockerComposeFile := path.Join(e.envProjPath, patterns.DockerComposeFile.Name)
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

		err = io.OverrideFile(path.Join(e.envProjPath, patterns.Makefile.Name), mkFile)
		if err != nil {
			return errors.Wrap(err, "error writing makefile")
		}
	}

	return nil
}
