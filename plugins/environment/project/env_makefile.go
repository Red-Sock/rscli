package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"
	"github.com/Red-Sock/rscli/plugins/environment/project/makefile"
)

type envMakefile struct {
	*makefile.Makefile
}

func (e *envMakefile) fetch(envProjPath string) (err error) {
	e.Makefile, err = makefile.ReadMakeFile(path.Join(envProjPath, envpatterns.Makefile.Name))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error getting makefile")
		}

		e.Makefile = makefile.MewEmptyMakefile()
	}

	return nil
}
