package scripts

import (
	"bytes"
	"fmt"
	"github.com/Red-Sock/rscli/internal/config"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrUnknownSource = errors.New("unknown client")

var skip = []byte("\n")

func RunSetUp(envs []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	subDir, dir := path.Split(wd)
	if dir == EnvDir {
		wd = subDir
	}
	m, err := validateEnv(wd, envs)
	if err != nil {
		return err
	}

	c, err := config.ReadConfig()
	if err != nil {
		return err
	}
	for _, name := range envs {
		err = addEnvironment(c, wd, name, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func addEnvironment(c *config.Config, wd, pName string, m map[int]string) error {
	envF, err := os.ReadFile(path.Join(wd, EnvDir, envExampleFile))
	if err != nil {
		return errors.Wrap(err, "can't open .env.example file")
	}

	envF = bytes.ReplaceAll(envF, []byte(projNameCapsPattern), []byte(strings.ToUpper(pName)))

	serverPort := []byte("SERVER_PORT")

	startIdx := bytes.Index(envF, serverPort) + len(serverPort) + 1
	endIdx := bytes.Index(envF[startIdx:], skip)
	if endIdx == -1 {
		endIdx = len(envF)
	} else {
		startIdx += endIdx
	}
	portBts := envF[startIdx:endIdx]
	port, err := strconv.Atoi(string(portBts))
	if err != nil {
		return err
	}

	envServerPortNameStart := bytes.LastIndex(envF[:startIdx], skip)
	if envServerPortNameStart == -1 {
		envServerPortNameStart = 0
	}

	envServerPortName := string(envF[envServerPortNameStart+1 : startIdx-1])

	for {
		_, ok := m[port]
		if !ok {
			m[port] = envServerPortName
			break
		}
		port++
	}

	envF = bytes.Replace(envF, append([]byte(envServerPortName+"="), portBts...), []byte(envServerPortName+"="+strconv.Itoa(port)), 1)
	{

		projectDir := path.Join(wd, pName, strings.ReplaceAll(c.Env.PathToClients, projNamePattern, pName))

		clients, err := getClients(projectDir)
		if err != nil {
			return err
		}

		for _, cl := range clients {
			envF = append(envF, skip...)
			port := cl.GetStartingPort()
			for {
				if _, ok := m[port]; !ok {
					m[port] = pName
					break
				}
				port++
			}

			envF = append(envF, []byte(fmt.Sprintf(cl.GetEnvs(), port))...)
		}
	}

	pathToEnv := path.Join(wd, EnvDir, pName, ".env")

	err = os.RemoveAll(pathToEnv)
	if err != nil {
		return err
	}

	err = os.WriteFile(pathToEnv, envF, 0755)
	if err != nil {
		return err
	}

	return nil
}

func validateEnv(wd string, envs []string) (map[int]string, error) {
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

func combineEnvFiles(wd string, envs []string) (string, error) {

	sb := &strings.Builder{}

	for _, dir := range envs {
		b, err := os.ReadFile(path.Join(wd, EnvDir, dir, ".env"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}
		sb.Write(b)
	}

	return sb.String(), nil
}

type Source interface {
	GetEnvs() string
	GetStartingPort() int
}

func getClients(pathToClients string) ([]Source, error) {

	dirs, err := os.ReadDir(pathToClients)
	if err != nil {
		return nil, err
	}
	clients := make([]Source, 0, len(dirs))

	for _, item := range dirs {
		if !item.IsDir() {
			continue
		}

		pattern, ok := clientPatterns[item.Name()]
		if !ok {
			return nil, ErrUnknownSource
		}

		clients = append(clients, pattern)
	}

	return clients, nil
}

func selectMakefile() []byte {
	if runtime.GOOS == "windows" {
		// TODO add windows support
		return nil
	} else {
		return makefile
	}
}

var clientPatterns = map[string]Source{
	"redis":    &redisEnvs{},
	"postgres": &pgEnvs{},
}
