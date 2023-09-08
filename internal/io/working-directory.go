package io

import (
	"os"
)

var workingDirectory string

func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic("cannot obtain working directory path:" + err.Error())
	}
}

func GetWd() string {
	return workingDirectory
}
