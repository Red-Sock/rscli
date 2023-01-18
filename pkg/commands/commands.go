package commands

import (
	"os"
	"path"
)

func init() {
	var err error
	rsCLI, err = os.Executable()
	if err != nil {
		panic(err)
	}
	_, rsCLI = path.Split(rsCLI)
}

func RsCLI() string {
	return rsCLI
}

var rsCLI string

const (
	FixUtil = "fix"
	GetUtil = "get"
	Delete  = "del"
)
