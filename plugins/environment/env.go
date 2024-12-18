package environment

import (
	"os"
	"path"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/compose"
	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/makefile"
)

const (
	PathFlag           = "path"
	ServiceInContainer = "service-enabled"
)

type GlobalEnvironment struct {
	io          io.IO
	rsCliConfig *config.RsCliConfig
	envDirPath  string

	composePatterns *compose.PatternManager
	environment     *env.Container
	makefile        *makefile.Makefile

	srcProjDirs []os.DirEntry
	envProjDirs []os.DirEntry
}

func NewGlobalEnv(io io.IO, cfg *config.RsCliConfig, envDirPath string) (*GlobalEnvironment, error) {
	c := &GlobalEnvironment{
		io:          io,
		rsCliConfig: cfg,
		envDirPath:  envDirPath,
	}

	return c, c.fetchSrcProjectDirs()
}

func (e *GlobalEnvironment) IsEnvExist() bool {
	for _, d := range e.srcProjDirs {
		if d.Name() == envpatterns.EnvDir {
			return true
		}
	}

	return false
}

func (e *GlobalEnvironment) fetchFiles() error {
	var err error

	err = e.fetchSrcProjectDirs()
	if err != nil {
		return rerrors.Wrap(err, "error fetching folders for environment")
	}

	err = e.fetchCompose()
	if err != nil {
		return rerrors.Wrap(err, "error fetching compose")
	}

	err = e.fetchDotEnv()
	if err != nil {
		return rerrors.Wrap(err, "error fetching dot env file")
	}

	err = e.fetchMakefile()
	if err != nil {
		return rerrors.Wrap(err, "error fetching makefile")
	}

	return nil

}

func (e *GlobalEnvironment) fetchSrcProjectDirs() (err error) {
	filter := func(dirs []os.DirEntry, srcProjDir string) ([]os.DirEntry, error) {
		var idx int
		for idx = 0; idx < len(dirs); idx++ {
			name := dirs[idx].Name()
			if dirs[idx].IsDir() && name != envpatterns.EnvDir {
				// validate whether this directory contains main file in specified location
				pathToMainFile := path.Join(srcProjDir, name,
					strings.ReplaceAll(e.rsCliConfig.Env.PathToMain, envpatterns.ProjNamePattern, name))

				fi, err := os.Stat(pathToMainFile)
				if err != nil {
					if !rerrors.Is(err, os.ErrNotExist) {
						return dirs, rerrors.Wrap(err, "error reading main.go file: "+pathToMainFile)
					}
				} else {
					if !fi.IsDir() {
						// definition of service to be added to projects dir list:
						// must be named NOT ComposePatterns.EnvDir
						// must be directory
						// must have proj_name/path_to_main
						continue
					}
				}
			}

			dirs[0], dirs[idx] = dirs[idx], dirs[0]
			dirs = dirs[1:]
			idx--
		}

		return dirs, nil
	}

	{
		e.srcProjDirs, err = os.ReadDir(path.Dir(e.envDirPath))
		if err != nil {
			return rerrors.Wrapf(err, "error reading directory projects %s ", e.envDirPath)
		}
		e.srcProjDirs, err = filter(e.srcProjDirs, path.Dir(e.envDirPath))
		if err != nil {
			return rerrors.Wrap(err, "error filtering source projects directories")
		}
	}

	{
		e.envProjDirs, err = os.ReadDir(e.envDirPath)
		if err != nil {
			if !rerrors.Is(err, os.ErrNotExist) {
				return rerrors.Wrapf(err, "error reading environment directory  %s ", e.envDirPath)
			}
		}

		e.envProjDirs, err = filter(e.envProjDirs, path.Dir(e.envDirPath))
		if err != nil {
			return rerrors.Wrap(err, "error filtering environment projects directories")
		}
	}

	return nil
}

func (e *GlobalEnvironment) fetchCompose() (err error) {
	composePatternsPath := path.Join(e.envDirPath, envpatterns.DockerComposeFile.Name)

	e.composePatterns, err = compose.ReadComposePatternsFromFile(composePatternsPath)
	if err != nil {
		return rerrors.Wrap(err, "error creating compose file at "+composePatternsPath)
	}

	return nil
}

func (e *GlobalEnvironment) fetchDotEnv() (err error) {
	builtIn, err := env.NewEnvContainer(envpatterns.EnvFile.Content)
	if err != nil {
		return rerrors.Wrap(err, "error parsing env container")
	}

	envPattern := path.Join(e.envDirPath, envpatterns.EnvFile.Name)
	globalEnv, err := env.ReadContainer(envPattern)
	if err != nil {
		return rerrors.Wrap(err, "can't open env file at "+envPattern)
	}

	for _, preDefined := range builtIn.GetContent() {
		for _, userDefined := range globalEnv.GetContent() {
			if preDefined.Name == userDefined.Name {
				builtIn.AppendRaw(preDefined.Name, userDefined.Value)
			}
		}
	}

	e.environment = builtIn

	return nil
}

func (e *GlobalEnvironment) fetchMakefile() (err error) {
	e.makefile, err = makefile.NewMakeFile(envpatterns.Makefile.Content)
	if err != nil {
		return rerrors.Wrap(err, "error parsing built in makefile")
	}

	userDefinedMakefilePath := path.Join(e.envDirPath, envpatterns.Makefile.Name)
	_, err = os.Stat(userDefinedMakefilePath)
	if err != nil {
		if rerrors.Is(err, os.ErrNotExist) {
			return nil
		}
		return rerrors.Wrap(err, "error getting stat on makefile")
	}

	m, err := makefile.ReadMakeFile(userDefinedMakefilePath)
	if err != nil {
		return rerrors.Wrap(err, "error parsing user defined config")
	}

	m.Merge(e.makefile)
	e.makefile = m

	return nil
}

func (e *GlobalEnvironment) getSpirits() []folder.Folder {
	return []folder.Folder{envpatterns.EnvFile, envpatterns.DockerComposeFile, envpatterns.Makefile}
}
