package io

import (
	"io/fs"
	"os"

	"go.redsock.ru/rerrors"
)

func CreateFileIfNotExists(pathToFile string, content []byte) error {
	fi, err := os.Stat(pathToFile)
	if err == nil {
		if fi.IsDir() {
			return rerrors.New(pathToFile + " already exists and it is folder")
		} else {
			return nil
		}
	}

	if !rerrors.Is(err, fs.ErrNotExist) {
		return rerrors.Wrap(err, "error reading file: "+pathToFile)
	}

	err = os.WriteFile(pathToFile, content, os.ModePerm)
	if err != nil {
		return rerrors.Wrap(err, "error writing Makefile")
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
			return rerrors.New(pth + " is not a directory")
		}
		return nil
	}

	if !rerrors.Is(err, fs.ErrNotExist) {
		return rerrors.Wrap(err, "error reading stat on path "+pth)
	}

	err = os.MkdirAll(pth, os.ModePerm)
	if err != nil {
		return rerrors.Wrap(err, "error making dir "+pth)
	}

	return nil
}
