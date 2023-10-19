package env

import (
	stderrs "errors"
	"os"
	"path"
	"sync"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/io"
)

func (c *Constructor) InitBasis() error {
	err := io.CreateFolderIfNotExists(c.envDirPath)

	for _, f := range c.getSpirits() {
		err = io.CreateFileIfNotExists(path.Join(c.envDirPath, f.Name), f.Content)
		if err != nil {
			return errors.Wrap(err, "error creating file "+f.Name+" file")
		}
	}

	return nil
}
func (c *Constructor) InitProjectsDirs() error {
	wg := &sync.WaitGroup{}
	errC := make(chan error, len(c.srcProjDirs))

	wg.Add(len(c.srcProjDirs))
	for _, d := range c.srcProjDirs {
		go func(d os.DirEntry) {
			defer wg.Done()

			err := c.initProjectDir(d)
			if err != nil {
				errC <- errors.Wrap(err, "error creating "+d.Name())
			}

		}(d)
	}

	wg.Wait()

	close(errC)

	var errs []error
	for errP := range errC {
		errs = append(errs, errP)
	}

	if len(errs) == 0 {
		return nil
	}

	return stderrs.Join(errors.New("error preparing projects dirs"), stderrs.Join(errs...))
}

func (c *Constructor) initProjectDir(d os.DirEntry) error {
	envProjDir := path.Join(c.envDirPath, d.Name())

	err := io.CreateFolderIfNotExists(envProjDir)
	if err != nil {
		return errors.Wrap(err, "error creating folder "+envProjDir)
	}

	var f []byte
	for _, spirit := range []patterns.File{patterns.Makefile} {

		f, err = os.ReadFile(path.Join(c.envDirPath, spirit.Name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return errors.Wrap(err, "error reading "+envProjDir+" file")
		}

		err = io.CreateFileIfNotExists(path.Join(envProjDir, spirit.Name), f)
		if err != nil {
			return errors.Wrap(err, "error reading "+envProjDir+" file")
		}
	}

	return nil
}
