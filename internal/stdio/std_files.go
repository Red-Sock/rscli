package stdio

import (
	"io/fs"
	"os"

	errors "github.com/Red-Sock/trace-errors"
)

func CreateFileIfNotExists(pathToFile string, content []byte) error {
	fi, err := os.Stat(pathToFile)
	if err == nil {
		if fi.IsDir() {
			return errors.New(pathToFile + " already exists and it is folder")
		} else {
			return nil
		}
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return errors.Wrap(err, "error reading file: "+pathToFile)
	}

	err = os.WriteFile(pathToFile, content, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing Makefile")
	}

	return nil
}

func OverrideFile(pth string, content []byte) error {
	err := os.RemoveAll(pth)
	if err != nil {
		return err
	}

	err = os.WriteFile(pth, content, 0755)
	if err != nil {
		return err
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
