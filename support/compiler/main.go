package main

import (
	"log"
	"os"
	"path"
	"strings"
)

// This supporting function compiles project pattern from
// pkg/service/project/pattern
// to pkg/service/project/pattern_c
// where all *.go files are renamed to *.go.pattern
// due to inability of go:embed to obtain project's *.go files

// !!! MUST BE EXECUTED FROM MAKEFILE FROM ROOT OF THIS PROJECT !!!

func main() {
	var oldPath, newPath string
	switch len(os.Args[1:]) {
	case 0:
		oldPath = findPatternFolder()
		newPath = oldPath + "_c"
	case 1:
		oldPath = os.Args[1]
	case 2:
		newPath = os.Args[2]
	}

	err := os.RemoveAll(newPath)
	if err != nil {
		log.Fatal(err)
	}

	movePattern(oldPath, newPath)
}

var pathToPattern = []string{
	"plugins", "project", "projpatterns", "pattern",
}

func findPatternFolder() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working directory", err)
	}

	for _, searchDir := range pathToPattern {
		dirs, err := os.ReadDir(wd)
		if err != nil {
			log.Fatal("error reading directory:", wd, err)
		}
		find := ""

		for _, d := range dirs {
			if !d.IsDir() {
				continue
			}
			n := d.Name()
			if n == searchDir {
				find = n
				break
			}
		}
		wd = path.Join(wd, find)
	}

	return wd
}

func movePattern(patternPath, newPath string) {
	err := os.MkdirAll(newPath, 0755)
	if err != nil {
		log.Fatal("error creating directory: ", newPath, err)
	}

	dirs, err := os.ReadDir(patternPath)
	if err != nil {
		log.Fatal("error reading directory", pathToPattern, err)
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
