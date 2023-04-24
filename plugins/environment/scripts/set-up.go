package scripts

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/config"
	_consts "github.com/Red-Sock/rscli/plugins/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/environment/scripts/patterns"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

var ErrUnknownSource = errors.New("unknown client")

var (
	lineSkip = []byte("\n")
)

type setupCommon struct {
	config                  *config.RsCliConfig
	workDir                 string
	envPattern              []byte
	composeFilePattern      []byte
	composeServicesPatterns map[string]patterns.ComposeService
	portToService           map[int]string
}

func RunSetUp(envs []string) (err error) {

	sc := setupCommon{
		workDir: wd,
	}
	subDir, dir := path.Split(wd)
	if dir == EnvDir {
		// if current dir is "environment" treat sub dir as WD
		sc.workDir = subDir
	}

	sc.composeServicesPatterns, err = patterns.NewComposeServices()
	if err != nil {
		return err
	}

	sc.portToService, err = getPortsFromEnv(wd, envs)
	if err != nil {
		return err
	}

	sc.config, err = config.ReadConfig(os.Args[1:])
	if err != nil {
		return err
	}

	sc.envPattern, err = os.ReadFile(path.Join(wd, EnvDir, envExampleFile))
	if err != nil {
		return errors.Wrapf(err, "can't open %s file", envExampleFile)
	}

	sc.composeFilePattern, err = os.ReadFile(path.Join(wd, EnvDir, composeExampleFile))
	if err != nil {
		return errors.Wrapf(err, "can't open %s file", envExampleFile)
	}

	for _, name := range envs {
		err = setUpEnv(name, sc)
		if err != nil {
			return err
		}
	}

	return nil
}

func setUpEnv(pName string, setup setupCommon) (err error) {
	var env *patterns.EnvService
	env, err = patterns.NewEnvService(envFile)
	if err != nil {
		return errors.Wrap(err, "error creating environment service")
	}

	var composeAssembler *patterns.ComposeAssembler
	composeAssembler, err = patterns.NewComposeAssembler(replaceProjectName(setup.composeFilePattern, pName), pName)
	if err != nil {
		return errors.Wrap(err, "error creating compose-file assembler")
	}

	projConf, err := pconfig.NewProjectConfig(path.Join(wd, pName, setup.config.Env.PathToConfig))
	if err != nil {
		return errors.Wrap(err, "error opening project configuration")
	}

	clients, err := getClients(setup.composeServicesPatterns, projConf)
	if err != nil {
		return errors.Wrap(err, "error assembling starting-compose-environment")
	}

	if clients == nil {
		return nil
	}

	for _, cl := range clients {

		composeEnvs := cl.GetEnvs().Content()

		for _, envVal := range composeEnvs {
			var p int
			p, err = strconv.Atoi(envVal.Value)
			if err == nil {
				for {
					if _, ok := setup.portToService[p]; !ok {
						setup.portToService[p] = pName
						env.Append(envVal.Name, strconv.Itoa(p))
						break
					}
					p++
				}
				continue
			}

			env.Append(envVal.Name, envVal.Value)
		}

		composeAssembler.AppendService(cl.Name, cl.GetCompose())
	}

	err = rewrite(replaceProjectName(env.MarshalEnv(), pName), path.Join(wd, EnvDir, pName, ".env"))
	if err != nil {
		return errors.Wrap(err, "error writing environment file")
	}

	bts, err := composeAssembler.Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshalling composer file")
	}

	bts = append([]byte("version: \"3.9\"\n\n"), bts...)
	err = rewrite(
		replaceProjectName(bts, pName),
		path.Join(wd, EnvDir, pName, dockerComposeFile),
	)
	if err != nil {
		return errors.Wrap(err, "error writing environment file")
	}

	return nil
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

func getClients(cs map[string]patterns.ComposeService, cfg interfaces.ProjectConfig) ([]patterns.ComposeService, error) {
	dsn, err := cfg.GetDataSourceOptions()
	if err != nil {
		return nil, err
	}

	clients := make([]patterns.ComposeService, 0, len(dsn))

	for _, item := range dsn {
		pattern, ok := cs[item.Type]
		if !ok {
			// TODO handle unknown data sources
			continue
		}
		pattern.Name = item.Name
		err = changeStartupParameters(&pattern, item)
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

func changeStartupParameters(pattern *patterns.ComposeService, conn pconfig.ConnectionOptions) error {
	switch conn.Type {
	case _consts.SourceNamePostgres:
		pattern.GetCompose().Environment[patterns.EnvironmentPostgresUser],
			pattern.GetCompose().Environment[patterns.EnvironmentPostgresPassword],
			_,
			_,
			pattern.GetCompose().Environment[patterns.EnvironmentPostgresDb] = _consts.ParsePgConnectionString(conn.ConnectionString)

	default:
		return ErrUnknownSource
	}
	return nil
}
