package project

import (
	"os"
	"path"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/makefile"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
)

type envMakefile struct {
	*makefile.Makefile
}

func (e *envMakefile) fetch(envProjPath string) (err error) {
	e.Makefile, err = makefile.ReadMakeFile(path.Join(envProjPath, patterns.Makefile.Name))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "error getting makefile")
		}

		e.Makefile = makefile.MewEmptyMakefile()
	}

	return nil
}
