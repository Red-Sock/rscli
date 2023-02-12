package scripts

import (
	"github.com/pkg/errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Red-Sock/rscli/internal/config"
)

var wd string

func init() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func ListProjects(pth string, cfg *config.Config) (projectNames []string, err error) {
	dirs, err := os.ReadDir(pth)
	if err != nil {
		return nil, err
	}

	projectNames = make([]string, 0, len(dirs))

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() && name != EnvDir {
			// validate whether this directory contains main file in specified location
			_, err = os.ReadFile(path.Join(pth, name, strings.ReplaceAll(cfg.Env.PathToMain, projNamePattern, name)))
			if err != nil {
				if err == os.ErrNotExist {
					continue
				}
				return nil, err
			}

			projectNames = append(projectNames, d.Name())
		}
	}

	return projectNames, nil
}

func ValidateEnvs(wd string, envs []string) error {
	_, err := getPortsFromEnv(wd, envs)
	return err
}

func getPortsFromEnv(wd string, envs []string) (map[int]string, error) {
	environment, err := combineEnvFiles(wd, envs)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(environment, "\n")
	m := make(map[int]string, len(rows))

	isValid := true

	var errs []error

	for _, row := range rows {
		if !strings.Contains(row, "=") {
			continue
		}

		splited := strings.Split(row, "=")

		if len(splited) != 2 {
			errs = append(errs, errors.New("ERROR PARSING .env file: "+splited[0]+" has no value"))
			isValid = false
			continue
		}
		if !strings.HasSuffix(splited[0], "PORT") {
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
			isValid = false
		}
		m[portInt] = splited[0]
	}

	if !isValid {
		os.Exit(1)
	}

	if errs == nil {
		return m, nil
	}
	err = errors.New("environment is invalid")

	for _, e := range errs {
		err = errors.Wrap(err, e.Error())
	}

	return nil, err
}
