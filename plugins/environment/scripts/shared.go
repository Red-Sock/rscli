package scripts

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Red-Sock/rscli/pkg/errors"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/environment/scripts/patterns"
)

const PortSuffix = "PORT"

var globalWD string

func init() {
	var err error
	globalWD, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func ListProjects(wd string, cfg *config.RsCliConfig) (projectNames []string, err error) {
	dirs, err := os.ReadDir(wd)
	if err != nil {
		return nil, errors.Wrap(err, "error reading dir")
	}

	projectNames = make([]string, 0, len(dirs))

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() && name != patterns.EnvDir {
			// validate whether this directory contains main file in specified location
			pathToMainFile := path.Join(wd, name, strings.ReplaceAll(cfg.Env.PathToMain, patterns.ProjNamePattern, name))
			_, err = os.ReadFile(pathToMainFile)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return nil, errors.Wrap(err, "error reading main.go file: "+pathToMainFile)
			}

			projectNames = append(projectNames, name)
		}
	}

	return projectNames, nil
}

func getPortsFromProjects(wd string, projectNames []string) (map[int]string, error) {
	environment, err := getEnvFilesCombined(wd, projectNames)
	if err != nil {
		return nil, errors.Wrap(err, "error obtaining combined env files")
	}

	rows := strings.Split(environment, "\n")
	m := make(map[int]string, len(rows))

	var errs []error

	for _, row := range rows {
		if !strings.Contains(row, "=") || strings.HasPrefix(row, "#") {
			continue
		}

		splited := strings.Split(row, "=")

		if len(splited) != 2 {
			errs = append(errs, errors.New("ERROR PARSING .env file: Row named "+splited[0]+" must have only 1 '=' symbol"))
			continue
		}

		if !strings.HasSuffix(splited[0], PortSuffix) {
			continue
		}

		portInt, err := strconv.Atoi(splited[1])
		if err != nil {
			errs = append(errs, errors.Wrap(err, "ERROR PARSING .env file port value "+splited[0]+" is not int but "+splited[1]))
			continue
		}
		if name, ok := m[portInt]; ok {
			errs = append(errs, errors.New("ERROR PARSING .env file. port : "+splited[1]+" is already assigned to "+name))
			continue
		}
		m[portInt] = splited[0]
	}

	if len(errs) == 0 {
		return m, nil
	}

	err = errors.New("environment is invalid")

	for _, e := range errs {
		err = errors.Wrap(err, e.Error())
	}

	return nil, err
}

func replaceProjectName(src []byte, newName string) []byte {
	b := make([]byte, len(src))
	copy(b, src)
	b = bytes.ReplaceAll(b, []byte(patterns.ProjNameCapsPattern), []byte(strings.ToUpper(newName)))
	b = bytes.ReplaceAll(b, []byte(patterns.ProjNamePattern), []byte(strings.ToLower(newName)))
	return b
}

func replaceProjectNameString(src string, newName string) string {
	src = strings.ReplaceAll(src, patterns.ProjNameCapsPattern, strings.ToUpper(newName))
	src = strings.ReplaceAll(src, patterns.ProjNamePattern, strings.ToLower(newName))
	return src
}

// TODO: make this function more secure and atomic.
// If WriteFile fails need to recover (e.g. create temp file and rename)
func rewrite(content []byte, pth string) (err error) {
	err = os.RemoveAll(pth)
	if err != nil {
		return err
	}

	err = os.WriteFile(pth, content, 0755)
	if err != nil {
		return err
	}

	return nil
}

func createFileIfNotExists(pathToFile string, content []byte) (err error) {
	_, err = os.ReadFile(pathToFile)
	if !errors.Is(err, fs.ErrNotExist) {
		return errors.Wrap(err, "error reading file: "+pathToFile)
	}

	err = os.WriteFile(pathToFile, content, 0755)
	if err != nil {
		return errors.Wrap(err, "error writing Makefile")
	}

	return nil

}
