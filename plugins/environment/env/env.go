package env

import (
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/environment/project/compose"
	"github.com/Red-Sock/rscli/plugins/environment/project/compose/env"
	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"
	"github.com/Red-Sock/rscli/plugins/environment/project/makefile"
)

const (
	PathFlag           = "path"
	ServiceInContainer = "service-enabled"
)

type Constructor struct {
	io  io.IO
	cfg *config.RsCliConfig

	composePatterns compose.PatternManager
	environment     *env.Container
	makefile        *makefile.Makefile

	envDirPath  string
	srcProjDirs []os.DirEntry
	EnvProjDirs []os.DirEntry
}

func NewConstructor(io io.IO, cfg *config.RsCliConfig) *Constructor {
	return &Constructor{
		io:  io,
		cfg: cfg,
	}
}

func (c *Constructor) FetchConstructor(cmd *cobra.Command, _ []string) error {
	var err error

	err = c.getWd(cmd)
	if err != nil {
		return errors.Wrap(err, "error fetching working directory")
	}

	err = c.fetchFolders()
	if err != nil {
		return errors.Wrap(err, "error fetching folders for environment")
	}

	err = c.fetchCompose()
	if err != nil {
		return errors.Wrap(err, "error fetching compose")
	}

	err = c.fetchDotEnv()
	if err != nil {
		return errors.Wrap(err, "error fetching dot env file")
	}

	err = c.fetchMakefile()
	if err != nil {
		return errors.Wrap(err, "error fetching makefile")
	}

	return nil

}

func (c *Constructor) IsEnvExist() bool {
	for _, d := range c.srcProjDirs {
		if d.Name() == envpatterns.EnvDir {
			return true
		}
	}

	return false
}

func (c *Constructor) getWd(cmd *cobra.Command) error {
	c.envDirPath = cmd.Flag(PathFlag).Value.String()

	if c.envDirPath == "" {
		c.envDirPath = io.GetWd()
	}

	if path.Base(c.envDirPath) != envpatterns.EnvDir {
		c.envDirPath = path.Join(c.envDirPath, envpatterns.EnvDir)
	}

	return nil
}

func (c *Constructor) fetchFolders() (err error) {
	filter := func(dirs []os.DirEntry, srcProjDir string) ([]os.DirEntry, error) {
		var idx int
		for idx = 0; idx < len(dirs); idx++ {
			name := dirs[idx].Name()
			if dirs[idx].IsDir() && name != envpatterns.EnvDir {
				// validate whether this directory contains main file in specified location
				pathToMainFile := path.Join(srcProjDir, name,
					strings.ReplaceAll(c.cfg.Env.PathToMain, envpatterns.ProjNamePattern, name))

				fi, err := os.Stat(pathToMainFile)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return dirs, errors.Wrap(err, "error reading main.go file: "+pathToMainFile)
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
		c.srcProjDirs, err = os.ReadDir(path.Dir(c.envDirPath))
		if err != nil {
			return errors.Wrapf(err, "error reading directory projects %s ", c.envDirPath)
		}
		c.srcProjDirs, err = filter(c.srcProjDirs, path.Dir(c.envDirPath))
		if err != nil {
			return errors.Wrap(err, "error filtering source projects directories")
		}
	}

	{
		c.EnvProjDirs, err = os.ReadDir(c.envDirPath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return errors.Wrapf(err, "error reading environment directory  %s ", c.envDirPath)
			}
		}

		c.EnvProjDirs, err = filter(c.EnvProjDirs, path.Dir(c.envDirPath))
		if err != nil {
			return errors.Wrap(err, "error filtering environment projects directories")
		}
	}

	return nil
}

func (c *Constructor) fetchCompose() (err error) {
	composePatternsPath := path.Join(c.envDirPath, envpatterns.DockerComposeFile.Name)

	c.composePatterns, err = compose.ReadComposePatternsFromFile(composePatternsPath)
	if err != nil {
		return errors.Wrap(err, "error creating compose file at "+composePatternsPath)
	}

	return nil
}

func (c *Constructor) fetchDotEnv() (err error) {
	builtIn, err := env.NewEnvContainer(envpatterns.EnvFile.Content)
	if err != nil {
		return errors.Wrap(err, "error parsing env container")
	}

	envPattern := path.Join(c.envDirPath, envpatterns.EnvFile.Name)
	globalEnv, err := env.ReadContainer(envPattern)
	if err != nil {
		return errors.Wrap(err, "can't open env file at "+envPattern)
	}

	for _, preDefined := range builtIn.GetContent() {
		for _, userDefined := range globalEnv.GetContent() {
			if preDefined.Name == userDefined.Name {
				builtIn.AppendRaw(preDefined.Name, userDefined.Value)
			}
		}
	}

	c.environment = builtIn

	return nil
}

func (c *Constructor) fetchMakefile() (err error) {
	c.makefile, err = makefile.NewMakeFile(envpatterns.Makefile.Content)
	if err != nil {
		return errors.Wrap(err, "error parsing built in makefile")
	}

	userDefinedMakefilePath := path.Join(c.envDirPath, envpatterns.Makefile.Name)
	_, err = os.Stat(userDefinedMakefilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.Wrap(err, "error getting stat on makefile")
	}

	m, err := makefile.ReadMakeFile(userDefinedMakefilePath)
	if err != nil {
		return errors.Wrap(err, "error parsing user defined config")
	}

	m.Merge(c.makefile)
	c.makefile = m

	return nil
}

func (c *Constructor) getSpirits() []folder.Folder {
	return []folder.Folder{envpatterns.EnvFile, envpatterns.DockerComposeFile, envpatterns.Makefile}
}
