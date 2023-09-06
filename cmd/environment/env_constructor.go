package environment

import (
	"os"
	"path"
	"runtime"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose"
	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/pkg/colors"
)

const (
	pathFlag = "path"
)

type envConstructor struct {
	cfg *config.RsCliConfig
	io  stdio.IO

	envDirPath  string
	srcProjDirs []os.DirEntry
	envProjDirs []os.DirEntry

	composePatterns compose.PatternManager
	envPatterns     *env.Container
}

func newEnvConstructor() *envConstructor {
	return &envConstructor{
		cfg: config.GetConfig(),
		io:  stdio.StdIO{},

		composePatterns: compose.PatternManager{},
		envPatterns:     &env.Container{},
	}
}

func (c *envConstructor) fetchConstructor(cmd *cobra.Command, _ []string) error {
	c.envDirPath = cmd.Flag(pathFlag).Value.String()

	if c.envDirPath == "" {
		c.envDirPath = stdio.GetWd()
	}

	if path.Base(c.envDirPath) != patterns.EnvDir {
		c.envDirPath = path.Join(c.envDirPath, patterns.EnvDir)
	}

	var err error

	c.srcProjDirs, err = os.ReadDir(path.Dir(c.envDirPath))
	if err != nil {
		return errors.Wrapf(err, "error reading directory projects %s ", c.envDirPath)
	}

	c.envProjDirs, err = os.ReadDir(c.envDirPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "error reading environment directory  %s ", c.envDirPath)
		}
	}

	composePatternsPath := path.Join(c.envDirPath, patterns.DockerComposeFile.Name)
	c.composePatterns, err = compose.ReadComposePatternsFromFile(composePatternsPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error reading compose file at "+composePatternsPath)
		}

	}

	envPattern := path.Join(c.envDirPath, patterns.EnvExampleFile)
	c.envPatterns, err = env.ReadContainer(envPattern)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "can't open env file at "+envPattern)
		}
	}

	err = c.filterFolders()
	if err != nil {
		return errors.Wrap(err, "error filtering folders")
	}

	return nil

}
func (c *envConstructor) filterFolders() error {
	filter := func(dirs []os.DirEntry, srcProjDir string) ([]os.DirEntry, error) {
		var idx int
		for idx = 0; idx < len(dirs); idx++ {
			name := dirs[idx].Name()
			if dirs[idx].IsDir() && name != patterns.EnvDir {
				// validate whether this directory contains main file in specified location
				pathToMainFile := path.Join(srcProjDir, name,
					strings.ReplaceAll(c.cfg.Env.PathToMain, patterns.ProjNamePattern, name))

				fi, err := os.Stat(pathToMainFile)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return dirs, errors.Wrap(err, "error reading main.go file: "+pathToMainFile)
					}
				} else {
					if !fi.IsDir() {
						// definition of service to be added to projects dir list:
						// must be named NOT composePatterns.EnvDir
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
	var err error

	c.srcProjDirs, err = filter(c.srcProjDirs, path.Dir(c.envDirPath))
	if err != nil {
		return errors.Wrap(err, "error filtering source projects directories")
	}

	c.envProjDirs, err = filter(c.envProjDirs, path.Dir(c.envDirPath))
	if err != nil {
		return errors.Wrap(err, "error filtering environment projects directories")
	}

	return nil
}

func (c *envConstructor) checkIfEnvExist() bool {
	for _, d := range c.srcProjDirs {
		if d.Name() == patterns.EnvDir {
			return true
		}
	}

	return false

}

func (c *envConstructor) askToRunTidy(cmd *cobra.Command, args []string, msg string) error {
	c.io.Println()
	c.io.PrintColored(colors.ColorYellow, msg+
		"!\nWant to run \"rscli env tidy\"? (Y)es/(N)o: ")

	for {
		resp, err := c.io.GetInput()
		if err != nil {
			return errors.Wrap(err, "error obtaining user input")
		}
		r := strings.ToLower(resp)[0]
		if r == 'y' {
			return c.runTidy(cmd, args)
		}

		if r == 'n' {
			return nil
		}
		c.io.PrintlnColored(colors.ColorRed, "No can't do "+resp+"! \"(Y)es\" or \"(N)o\":")
	}

	return nil
}

func (c *envConstructor) selectMakefile() patterns.File {
	if runtime.GOOS == "windows" {
		// TODO add windows support
		return patterns.Makefile
	} else {
		return patterns.Makefile
	}
}

func (c *envConstructor) getSpirits() []patterns.File {
	return []patterns.File{patterns.EnvFile, patterns.DockerComposeFile, c.selectMakefile()}
}
