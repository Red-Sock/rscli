package main

import (
	"log"
	"os"
	"path"
	"strings"
)

// This supporting function compiles project pattern from
// https://github.com/Red-Sock/project-pattern/tree/go
// to pkg/service/project/pattern_c
// where all *.go files are renamed to *.go.pattern
// due to inability of go:embed to obtain project's *.go files

// !!! MUST BE EXECUTED FROM MAKEFILE IN ROOT OF PROJECT !!!

func main() {
	var oldPath, newPath string
	if len(os.Args[1:]) < 2 {
		println("path to project pattern and compiling folder has to be specified")
		os.Exit(1)
	}
	oldPath = os.Args[1]
	newPath = os.Args[2]

	err := os.RemoveAll(newPath)
	if err != nil {
		log.Fatal(err)
	}

	movePattern(oldPath, newPath)
}

func movePattern(patternPath, newPath string) {
	err := os.MkdirAll(newPath, 0755)
	if err != nil {
		log.Fatal("error creating directory: ", newPath, err)
	}

	dirs, err := os.ReadDir(patternPath)
	if err != nil {
		log.Fatal("error reading directory", patternPath, err)
	}

	for _, d := range dirs {
		itemName := d.Name()
		newItemPath := path.Join(newPath, itemName)

		if d.IsDir() {
			movePattern(path.Join(patternPath, itemName), newItemPath)
		} else {
			var b []byte
			pathToFile := path.Join(patternPath, itemName)
			b, err = os.ReadFile(pathToFile)
			if err != nil {
				log.Fatal("error reading file:", pathToFile, err)
			}

			if hasOneOfSuffixes(itemName, ".go", "go.mod") {
				newItemPath += ".pattern"
			}

			err = os.WriteFile(newItemPath, b, 0755)
			if err != nil {
				log.Fatal("error writing file: ", newItemPath, err)
			}
		}
	}

}

func hasOneOfSuffixes(in string, sufs ...string) bool {
	for _, s := range sufs {
		if strings.HasSuffix(in, s) {
			return true
		}
	}
	return false
}
