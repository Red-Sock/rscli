package scripts

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/Red-Sock/trace-errors"

	patterns3 "github.com/Red-Sock/rscli/cmd/environment/patterns"
	"github.com/Red-Sock/rscli/internal/config"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
	patterns2 "github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var ErrUnknownSource = errors.New("unknown client")

// environmentSetupConfig - configuration info needed for setting up an environment
type environmentSetupConfig struct {
	config                  *config.RsCliConfig
	workDir                 string
	envPattern              []byte
	composeServicesPatterns map[string]patterns3.ComposePatterns
	portToService           portManager
}

// RunSetUp - runs actual environment setup.
// Prepares datasource dependencies files - .env docker-compose.yaml and etc.
func RunSetUp(projectNames []string) (err error) {
	var sc = environmentSetupConfig{
		workDir: globalWD,
	}

	sc.config = config.GetConfig()

	subDir, dir := path.Split(sc.workDir)
	if dir == patterns3.EnvDir {
		// If current working directory is called "environment" -> treat parent directory as WD
		sc.workDir = subDir
	}

	err = createEnvFolders(sc.workDir, sc.config)
	if err != nil {
		return errors.Wrap(err, "error creating env folders")
	}

	// Docker-compose patterns
	sc.composeServicesPatterns, err = patterns3.NewComposePatterns(sc.workDir)
	if err != nil {
		return errors.Wrap(err, "error creating compose patterns")
	}

	// .env patterns
	sc.envPattern, err = os.ReadFile(path.Join(sc.workDir, patterns3.EnvDir, patterns3.EnvExampleFile))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "can't open %s file", patterns3.EnvExampleFile)
	}

	sc.portToService, err = getPortsFromProjects(sc.workDir, projectNames)
	if err != nil {
		return err
	}

	for _, projName := range projectNames {
		err = setUpEnvForProject(projName, sc)
		if err != nil {
			return err
		}
	}

	return nil
}

func setUpEnvForProject(pName string, setup environmentSetupConfig) (err error) {
	var envAssembler *patterns3.EnvService
	envAssembler, err = patterns3.NewEnvService(patterns3.EnvFile.Content)
	if err != nil {
		return errors.Wrap(err, "error creating environment service")
	}

	var composeAssembler *patterns3.ComposeAssembler
	composeAssembler, err = patterns3.NewComposeAssembler(replaceProjectName(patterns3.DockerComposeFile.Content, pName), pName)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	configPath := path.Join(setup.workDir, pName, setup.config.Env.PathToConfig)
	projConf, err := pconfig.ParseConfig(configPath)
	if err != nil {
		return errors.Wrap(err, "error parsing users's config")
	}

	{
		clients, err := getClients(setup.composeServicesPatterns, projConf)
		if err != nil {
			return errors.Wrap(err, "error assembling starting-compose-environment")
		}

		for _, resource := range clients {

			composeEnvs := resource.GetEnvs().Content()

			for _, envRow := range composeEnvs {
				if strings.HasSuffix(envRow.Name, PortSuffix) {
					var p int
					p, err = strconv.Atoi(envRow.Value)
					if err != nil {
						return errors.Wrap(err, "error parsing .env file: port value for "+envRow.Name+" must be int but it is "+envRow.Value)
					}

					p = setup.portToService.GetNextPort(p, pName)

					envAssembler.Append(envRow.Name, strconv.Itoa(p))
				}

				envAssembler.Append(envRow.Name, envRow.Value)
			}

			{
				composeAssembler.AppendService(resource.Name, resource.GetCompose())
			}
		}
	}

	{
		opts := projConf.GetServerOptions()

		for _, srvOpt := range opts {
			if srvOpt.Port == 0 {
				continue
			}
			portName := patterns3.ProjNameCapsPattern + "_" + strings.ToUpper(srvOpt.Name) + "_" + PortSuffix
			composeAssembler.Services[pName].Ports = append(composeAssembler.Services[pName].Ports, addEnvironmentBrackets(portName)+":"+strconv.Itoa(int(srvOpt.Port)))
			envAssembler.Append(portName, strconv.Itoa(int(srvOpt.Port)))
		}
	}

	pathToProjectEnvFile := path.Join(setup.workDir, patterns3.EnvDir, pName, patterns3.EnvFile.Name)
	err = rewrite(replaceProjectName(envAssembler.MarshalEnv(), pName), pathToProjectEnvFile)
	if err != nil {
		return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
	}

	composeFile, err := composeAssembler.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	pathToDockerComposeFile := path.Join(setup.workDir, patterns3.EnvDir, pName, patterns3.DockerComposeFile.Name)
	err = rewrite(replaceProjectName(composeFile, pName), pathToDockerComposeFile)
	if err != nil {
		return errors.Wrap(err, "error writing docker compose file file")
	}

	return nil
}

func getEnvFilesCombined(wd string, projectNames []string) (string, error) {
	sb := &strings.Builder{}

	for _, dir := range projectNames {
		projEnvFile, err := os.ReadFile(path.Join(wd, patterns3.EnvDir, dir, patterns3.EnvFile.Name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", errors.Wrap(err, "error reading environment file")
		}

		sb.Write(projEnvFile)
		if projEnvFile[len(projEnvFile)-1] != '\n' {
			sb.WriteByte('\n')
		}
	}

	return sb.String(), nil
}

func getClients(composePatterns map[string]patterns3.ComposePatterns, cfg *pconfig.Config) ([]patterns3.ComposePatterns, error) {
	dsns, err := cfg.GetDataSourceOptions()
	if err != nil {
		return nil, err
	}

	clients := make([]patterns3.ComposePatterns, 0, len(dsns))

	for _, dsn := range dsns {
		pattern, ok := composePatterns[dsn.Type]
		if !ok {
			// TODO handle unknown data sources
			continue
		}
		pattern.Name = dsn.Name
		err = insertEnvironmentValues(&pattern, dsn)
		if err != nil {
			return nil, errors.Wrap(err, "error inserting environment values to compose")
		}

		clients = append(clients, pattern)
	}

	for _, c := range clients {
		c.GetEnvs().RemoveEmpty()
	}

	return clients, nil
}

func insertEnvironmentValues(pattern *patterns3.ComposePatterns, conn pconfig.ConnectionOptions) error {
	switch conn.Type {
	case patterns2.SourceNamePostgres:
		user, pwd, _, _, dbName := pconfig.ParsePgConnectionString(conn.ConnectionString)
		env := pattern.GetEnvs()
		composeEnv := pattern.GetCompose().Environment

		const dsPgName = "DS_POSTGRES_NAME_CAPS"
		{
			varName := composeEnv[patterns3.EnvironmentPostgresUser]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns3.EnvironmentPostgresUser] = varName

			env.Append(removeEnvironmentBrackets(varName), user)
		}
		{
			varName := composeEnv[patterns3.EnvironmentPostgresPassword]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns3.EnvironmentPostgresPassword] = varName

			env.Append(removeEnvironmentBrackets(varName), pwd)
		}
		{
			varName := composeEnv[patterns3.EnvironmentPostgresDb]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns3.EnvironmentPostgresDb] = varName

			env.Append(removeEnvironmentBrackets(varName), dbName)
		}
	default:
		return ErrUnknownSource
	}
	return nil
}

func removeEnvironmentBrackets(in string) string {
	if in[:2] == "${" && in[len(in)-1] == '}' {
		return in[2 : len(in)-1]
	}

	return in
}
func addEnvironmentBrackets(in string) string {
	return "${" + in + "}"
}

func selectMakefile() patterns3.File {
	if runtime.GOOS == "windows" {
		// TODO add windows support
		return patterns3.Makefile
	} else {
		return patterns3.Makefile
	}
}

type portManager map[int]string

func (p portManager) GetNextPort(in int, projName string) int {
	for {
		// if such port already exists - increment it
		if _, ok := p[in]; !ok {
			p[in] = projName
			return in
		}
		in++
	}
}
