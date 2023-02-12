package scripts

import (
	"bytes"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/environment/scripts/patterns"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrUnknownSource = errors.New("unknown client")

var lineSkip = []byte("\n")

func RunSetUp(envs []string) error {
	cs, err := patterns.NewComposeService()
	if err != nil {
		return err
	}

	wd := wd

	subDir, dir := path.Split(wd)
	if dir == EnvDir {
		// if current dir is "environment" treat sub dir as WD
		wd = subDir
	}

	m, err := getPortsFromEnv(wd, envs)
	if err != nil {
		return err
	}

	c, err := config.ReadConfig()
	if err != nil {
		return err
	}

	for _, name := range envs {
		err = addEnvironment(c, cs, wd, name, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func addEnvironment(c *config.Config, cs map[string]patterns.ComposeService, wd, pName string, portToService map[int]string) error {
	envF, err := os.ReadFile(path.Join(wd, EnvDir, envExampleFile))
	if err != nil {
		return errors.Wrap(err, "can't open .env.example file")
	}

	envF = bytes.ReplaceAll(envF, []byte(projNameCapsPattern), []byte(strings.ToUpper(pName)))

	serverPortString := []byte("SERVER_PORT")

	startIdx := bytes.Index(envF, serverPortString) + len(serverPortString) + 1
	endIdx := bytes.Index(envF[startIdx:], lineSkip)
	if endIdx == -1 {
		endIdx = len(envF)
	} else {
		startIdx += endIdx
	}
	portBts := envF[startIdx:endIdx]
	serverPort, err := strconv.Atoi(string(portBts))
	if err != nil {
		return err
	}

	envServerPortNameStart := bytes.LastIndex(envF[:startIdx], lineSkip)
	if envServerPortNameStart == -1 {
		envServerPortNameStart = 0
	}

	envServerPortName := string(envF[envServerPortNameStart+1 : startIdx-1])

	for {
		_, ok := portToService[serverPort]
		if !ok {
			portToService[serverPort] = envServerPortName
			break
		}
		serverPort++
	}

	envF = bytes.Replace(envF, append([]byte(envServerPortName+"="), portBts...), []byte(envServerPortName+"="+strconv.Itoa(serverPort)), 1)
	{

		projectDir := path.Join(wd, pName, strings.ReplaceAll(c.Env.PathToClients, projNamePattern, pName))

		clients, err := getClients(cs, projectDir)
		if err != nil {
			return err
		}

		// todo
		for _, cl := range clients {
			envF = append(envF, lineSkip...)

			envsStrings := strings.Split(cl.GetEnvs(), "\n")
			ports := cl.GetStartingPorts()

			for idx, p := range ports {
				for {
					if _, ok := portToService[p]; !ok {
						portToService[p] = pName
						envF = append(envF, []byte(envsStrings[idx]+"="+strconv.Itoa(p))...)
						break
					}
					p++
				}

			}
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

func addCompose() {

}

func combineEnvFiles(wd string, envs []string) (string, error) {
	sb := &strings.Builder{}

	for _, dir := range envs {
		projEnvFile, err := os.ReadFile(path.Join(wd, EnvDir, dir, EnvFile))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", errors.Wrap(err, "error reading environment file")
		}

		sb.Write(projEnvFile)
	}

	return sb.String(), nil
}

func getClients(cs map[string]patterns.ComposeService, pathToClients string) ([]patterns.ComposeService, error) {

	dirs, err := os.ReadDir(pathToClients)
	if err != nil {
		return nil, err
	}

	clients := make([]patterns.ComposeService, 0, len(dirs))

	for _, item := range dirs {
		if !item.IsDir() {
			continue
		}

		pattern, ok := cs[item.Name()]
		if !ok {
			continue
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
