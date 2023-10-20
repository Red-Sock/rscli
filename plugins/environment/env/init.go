package env

import (
	stderrs "errors"
	"os"
	"path"
	"sync"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"

	"github.com/Red-Sock/rscli/internal/io"
)

func (e *GlobalEnvironment) Init() error {
	err := e.initBasis()
	if err != nil {
		return errors.Wrap(err, "error initiating basis")
	}

	err = e.initProjectsDirs()
	if err != nil {
		return errors.Wrap(err, "error initiating project dirs")
	}

	return nil
}

func (e *GlobalEnvironment) initBasis() error {
	err := io.CreateFolderIfNotExists(e.envDirPath)

	for _, f := range e.getSpirits() {
		err = io.CreateFileIfNotExists(path.Join(e.envDirPath, f.Name), f.Content)
		if err != nil {
			return errors.Wrap(err, "error creating file "+f.Name+" file")
		}
	}

	return nil
}
func (e *GlobalEnvironment) initProjectsDirs() error {
	wg := &sync.WaitGroup{}
	errC := make(chan error, len(e.srcProjDirs))

	wg.Add(len(e.srcProjDirs))
	for _, d := range e.srcProjDirs {
		go func(d os.DirEntry) {
			defer wg.Done()

			err := e.initProjectDir(d)
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

func (e *GlobalEnvironment) initProjectDir(d os.DirEntry) error {
	envProjDir := path.Join(e.envDirPath, d.Name())

	err := io.CreateFolderIfNotExists(envProjDir)
	if err != nil {
		return errors.Wrap(err, "error creating folder "+envProjDir)
	}

	var f []byte
	for _, spirit := range []folder.Folder{envpatterns.Makefile} {

		f, err = os.ReadFile(path.Join(e.envDirPath, spirit.Name))
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
