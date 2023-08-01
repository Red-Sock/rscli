package scripts

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/plugins/environment/scripts/patterns"
	"github.com/Red-Sock/rscli/plugins/project/config/pkg/configstructs"
	"github.com/Red-Sock/rscli/plugins/project/config/pkg/const"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

var ErrUnknownSource = errors.New("unknown client")

// environmentSetupConfig - configuration info needed for setting up an environment
type environmentSetupConfig struct {
	config                  *config.RsCliConfig
	workDir                 string
	envPattern              []byte
	composeServicesPatterns map[string]patterns.ComposePatterns
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
	if dir == patterns.EnvDir {
		// If current working directory is called "environment" -> treat parent directory as WD
		sc.workDir = subDir
	}

	err = createEnvFolders(sc.workDir, sc.config)
	if err != nil {
		return errors.Wrap(err, "error creating env folders")
	}

	// Docker-compose patterns
	sc.composeServicesPatterns, err = patterns.NewComposePatterns(sc.workDir)
	if err != nil {
		return errors.Wrap(err, "error creating compose patterns")
	}

	// .env patterns
	sc.envPattern, err = os.ReadFile(path.Join(sc.workDir, patterns.EnvDir, patterns.EnvExampleFile))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "can't open %s file", patterns.EnvExampleFile)
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
	var envAssembler *patterns.EnvService
	envAssembler, err = patterns.NewEnvService(patterns.EnvFile.Content)
	if err != nil {
		return errors.Wrap(err, "error creating environment service")
	}

	var composeAssembler *patterns.ComposeAssembler
	composeAssembler, err = patterns.NewComposeAssembler(replaceProjectName(patterns.DockerComposeFile.Content, pName), pName)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	projConf, err := pconfig.NewProjectConfig(path.Join(setup.workDir, pName, setup.config.Env.PathToConfig))
	if err != nil {
		return errors.Wrap(err, "error opening project configuration")
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
		opts, err := projConf.GetServerOptions()
		if err != nil {
			return errors.Wrap(err, "error obtaining server options")
		}
		for _, srvOpt := range opts {
			if srvOpt.Port == 0 {
				continue
			}
			portName := patterns.ProjNameCapsPattern + "_" + strings.ToUpper(srvOpt.Name) + "_" + PortSuffix
			composeAssembler.Services[pName].Ports = append(composeAssembler.Services[pName].Ports, addEnvironmentBrackets(portName)+":"+strconv.Itoa(int(srvOpt.Port)))
			envAssembler.Append(portName, strconv.Itoa(int(srvOpt.Port)))
		}
	}

	pathToProjectEnvFile := path.Join(setup.workDir, patterns.EnvDir, pName, patterns.EnvFile.Name)
	err = rewrite(replaceProjectName(envAssembler.MarshalEnv(), pName), pathToProjectEnvFile)
	if err != nil {
		return errors.Wrap(err, "error writing environment file: "+pathToProjectEnvFile)
	}

	composeFile, err := composeAssembler.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	pathToDockerComposeFile := path.Join(setup.workDir, patterns.EnvDir, pName, patterns.DockerComposeFile.Name)
	err = rewrite(replaceProjectName(composeFile, pName), pathToDockerComposeFile)
	if err != nil {
		return errors.Wrap(err, "error writing docker compose file file")
	}

	return nil
}

func getEnvFilesCombined(wd string, projectNames []string) (string, error) {
	sb := &strings.Builder{}

	for _, dir := range projectNames {
		projEnvFile, err := os.ReadFile(path.Join(wd, patterns.EnvDir, dir, patterns.EnvFile.Name))
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

func getClients(composePatterns map[string]patterns.ComposePatterns, cfg interfaces.ProjectConfig) ([]patterns.ComposePatterns, error) {
	dsns, err := cfg.GetDataSourceOptions()
	if err != nil {
		return nil, err
	}

	clients := make([]patterns.ComposePatterns, 0, len(dsns))

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

func insertEnvironmentValues(pattern *patterns.ComposePatterns, conn configstructs.ConnectionOptions) error {
	switch conn.Type {
	case _const.SourceNamePostgres:
		user, pwd, _, _, dbName := _const.ParsePgConnectionString(conn.ConnectionString)
		env := pattern.GetEnvs()
		composeEnv := pattern.GetCompose().Environment

		const dsPgName = "DS_POSTGRES_NAME_CAPS"
		{
			varName := composeEnv[patterns.EnvironmentPostgresUser]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvironmentPostgresUser] = varName

			env.Append(removeEnvironmentBrackets(varName), user)
		}
		{
			varName := composeEnv[patterns.EnvironmentPostgresPassword]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvironmentPostgresPassword] = varName

			env.Append(removeEnvironmentBrackets(varName), pwd)
		}
		{
			varName := composeEnv[patterns.EnvironmentPostgresDb]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvironmentPostgresDb] = varName

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

func selectMakefile() patterns.File {
	if runtime.GOOS == "windows" {
		// TODO add windows support
		return patterns.Makefile
	} else {
		return patterns.Makefile
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
