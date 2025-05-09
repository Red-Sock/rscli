package project

import (
	"path"

	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/makefile"
	"github.com/Red-Sock/rscli/internal/ports"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
)

var ErrNoConfig = rerrors.New("no config found")

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
			Network:  globalComposePatternManager.Network,
		},
		rscliConfig:                 rscliConfig,
		globalPortManager:           globalPortManager,
		globalMakefile:              globalMakefile,
		globalComposePatternManager: globalComposePatternManager,
	}

	err = p.Environment.fetch(globalEnv, pathToProjectEnv)
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching .env file")
	}

	err = p.Config.fetch(rscliConfig, pathToProjectEnv, pathToProject)
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching config")
	}

	err = p.Makefile.fetch(pathToProjectEnv)
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching makefile")
	}

	return p, nil
}

func (e *ProjEnv) Tidy(serviceEnabled bool) error {
	projName := path.Base(e.pathToProjInEnv)

	err := e.tidyResources(serviceEnabled)
	if err != nil {
		return rerrors.Wrap(err, "error doing tidy on resources")
	}

	err = e.tidyServerAPIs()
	if err != nil {
		return rerrors.Wrap(err, "error doing tidy on server api")
	}

	err = e.flush(projName)
	if err != nil {
		return rerrors.Wrap(err, "error flushing files")
	}

	err = e.tidyMigrationDirs()
	if err != nil {
		return rerrors.Wrap(err, "error")
	}

	return nil
}

func (e *ProjEnv) flush(projName string) (err error) {
	{
		pathToProjectEnvFile := path.Join(e.pathToProjInEnv, envpatterns.EnvFile.Name)

		envBytes := e.Environment.MarshalEnv()
		if len(envBytes) != 0 {
			err = io.OverrideFile(pathToProjectEnvFile, renamer.ReplaceProjectNameShort(envBytes, projName))
			if err != nil {
				return rerrors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
			}
		}
	}

	{
		var composeFile []byte
		composeFile, err = e.Compose.Marshal()
		if err != nil {
			return rerrors.Wrap(err, "error marshalling composer file")
		}

		pathToDockerComposeFile := path.Join(e.pathToProjInEnv, envpatterns.DockerComposeFile.Name)
		err = io.OverrideFile(pathToDockerComposeFile, renamer.ReplaceProjectNameShort(composeFile, projName))
		if err != nil {
			return rerrors.Wrap(err, "error writing docker compose file file")
		}
	}

	{
		var b []byte
		b, err = e.Config.Marshal()
		if err != nil {
			return rerrors.Wrap(err, "error marshalling env config")
		}

		err = io.OverrideFile(e.Config.pth, b)
		if err != nil {
			return rerrors.Wrap(err, "error writing env config")
		}
	}

	{
		e.tidyMakeFile()

		var mkFile []byte
		mkFile, err = e.Makefile.Marshal()
		if err != nil {
			return rerrors.Wrap(err, "error marshalling makefile")
		}

		mkFile = renamer.ReplaceProjectNameShort(mkFile, projName)

		err = io.OverrideFile(path.Join(e.pathToProjInEnv, envpatterns.Makefile.Name), mkFile)
		if err != nil {
			return rerrors.Wrap(err, "error writing makefile")
		}
	}

	return nil
}
