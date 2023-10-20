package project

import (
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/makefile"
	"github.com/Red-Sock/rscli/internal/ports"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

var ErrNoConfig = errors.New("no config found")

type ProjEnv struct {
	projName        string
	pathToProjInEnv string
	pathToProjSrc   string

	Compose     *compose.Compose
	Environment envVariables
	Makefile    envMakefile
	Config      envConfig

	rscliConfig                 *config.RsCliConfig
	globalComposePatternManager *compose.PatternManager
	globalMakefile              *makefile.Makefile

	globalPortManager *ports.PortManager
}

func LoadProjectEnvironment(
	rscliConfig *config.RsCliConfig,
	globalEnv *env.Container,
	globalMakefile *makefile.Makefile,
	globalComposePatternManager *compose.PatternManager,

	globalPortManager *ports.PortManager,

	pathToProjectEnv string,
	pathToProject string,

) (p *ProjEnv, err error) {
	p = &ProjEnv{
		projName:        path.Base(pathToProject),
		pathToProjInEnv: pathToProjectEnv,
		pathToProjSrc:   pathToProject,

		Compose: &compose.Compose{
			Services: map[string]*compose.ContainerSettings{},
			Network:  map[string]interface{}{},
		},
		rscliConfig:                 rscliConfig,
		globalPortManager:           globalPortManager,
		globalMakefile:              globalMakefile,
		globalComposePatternManager: globalComposePatternManager,
	}

	err = p.Environment.fetch(globalEnv, pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching .env file")
	}

	err = p.Config.fetch(rscliConfig, pathToProjectEnv, pathToProject)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching config")
	}

	err = p.Makefile.fetch(pathToProjectEnv)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching makefile")
	}

	return p, nil
}

func (e *ProjEnv) Tidy(serviceEnabled bool) error {
	projName := path.Base(e.pathToProjInEnv)

	err := e.tidyResources(serviceEnabled)
	if err != nil {
		return errors.Wrap(err, "error doing tidy on resources")
	}

	err = e.tidyServerAPIs()
	if err != nil {
		return errors.Wrap(err, "error doing tidy on server api")
	}

	err = e.tidyService()
	if err != nil {
		return errors.Wrap(err, "error doing tidy on service")
	}

	err = e.flush(projName)

	return nil
}

func (e *ProjEnv) flush(projName string) (err error) {
	{
		pathToProjectEnvFile := path.Join(e.pathToProjInEnv, envpatterns.EnvFile.Name)

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

		pathToDockerComposeFile := path.Join(e.pathToProjInEnv, envpatterns.DockerComposeFile.Name)
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
		e.tidyMakeFile()

		var mkFile []byte
		mkFile, err = e.Makefile.Marshal()
		if err != nil {
			return errors.Wrap(err, "error marshalling makefile")
		}

		err = io.OverrideFile(path.Join(e.pathToProjInEnv, envpatterns.Makefile.Name), mkFile)
		if err != nil {
			return errors.Wrap(err, "error writing makefile")
		}
	}

	return nil
}
