package stdio

import (
	"io/fs"
	"os"

	errors "github.com/Red-Sock/trace-errors"
)

func CreateFileIfNotExists(pathToFile string, content []byte) error {
	fi, err := os.Stat(pathToFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Wrap(err, "error reading file: "+pathToFile)
		}
	} else {
		if fi.IsDir() {
			return errors.New(pathToFile + " already exists and it is folder")
		}
	}

	err = os.WriteFile(pathToFile, content, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing Makefile")
	}

	return nil
}

func CreateFolderIfNotExists(pth string) error {
	fi, err := os.Stat(pth)
	if err == nil {
		if !fi.IsDir() {
			return errors.New(pth + " is not a directory")
		}
		return nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return errors.Wrap(err, "error reading stat on path "+pth)
	}

	err = os.MkdirAll(pth, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error making dir "+pth)
	}

	return nil
}
