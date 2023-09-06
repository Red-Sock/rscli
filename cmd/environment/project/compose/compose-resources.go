package compose

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/utils/nums"
	pconfig "github.com/Red-Sock/rscli/plugins/project/processor/config"
	projPatterns "github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

const (
	servicesPart = "services"
)

var (
	ErrInvalidComposeFileFormat = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposeEnvFormat  = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
	ErrUnknownSource            = errors.New("unknown client")
)

type PatternManager struct {
	Patterns map[string]Pattern
}

type Pattern struct {
	Name    string
	content ContainerSettings
	envs    *env.Container
}

func ReadComposePatternsFromFile(pth string) (out PatternManager, err error) {
	out.Patterns = map[string]Pattern{}

	// Basic compose examples: rscli built-in
	out.Patterns, err = extractComposePatternsFromFile(patterns.BuildInComposeExamples)
	if err != nil {
		return PatternManager{}, errors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	// User's defined compose examples from
	userDockerComposeExample, err := os.ReadFile(pth)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}

		return out, errors.Wrap(err, "error reading user defined compose file")
	}

	userServicesDefinitions, err := extractComposePatternsFromFile(userDockerComposeExample)
	if err != nil {
		return out, errors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	for serviceName, content := range userServicesDefinitions {
		// override predefined service description with user's defined one
		out.Patterns[serviceName] = content
	}

	return out, nil
}

func (p *Pattern) GetEnvs() *env.Container {
	return p.envs
}

func (p *Pattern) GetCompose() ContainerSettings {
	return p.content
}

func (p *Pattern) insertEnvironmentValues(conn pconfig.ConnectionOptions) error {
	switch conn.Type {
	case projPatterns.SourceNamePostgres:
		user, pwd, _, _, dbName := pconfig.ParsePgConnectionString(conn.ConnectionString)
		env := p.GetEnvs()
		composeEnv := p.GetCompose().Environment

		// TODO подумать и двинуть это как-нибудь в сторону структуры
		const dsPgName = "DS_POSTGRES_NAME_CAPS"
		{
			varName := composeEnv[patterns.EnvVarPostgresUser]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvVarPostgresUser] = varName

			env.Append(removeEnvironmentBrackets(varName), user)
		}
		{
			varName := composeEnv[patterns.EnvVarPostgresPassword]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvVarPostgresPassword] = varName

			env.Append(removeEnvironmentBrackets(varName), pwd)
		}
		{
			varName := composeEnv[patterns.EnvVarPostgresDB]
			varName = strings.ToUpper(strings.ReplaceAll(varName, dsPgName, dbName))
			composeEnv[patterns.EnvVarPostgresDB] = varName

			env.Append(removeEnvironmentBrackets(varName), dbName)
		}
	default:
		return ErrUnknownSource
	}
	return nil
}

func (c *PatternManager) GetServiceDependencies(cfg *pconfig.Config) ([]Pattern, error) {
	// TODO переделать на вариант без ТИПА, нужен интерфес -> имя, структура конфига подключения, порты
	dsns, err := cfg.GetDataSourceOptions()
	if err != nil {
		return nil, err
	}

	clients := make([]Pattern, 0, len(dsns))

	for _, dsn := range dsns {
		pattern, ok := c.Patterns[dsn.Type]
		if !ok {
			// TODO handle unknown data sources
			continue
		}
		pattern.Name = dsn.Name
		err = pattern.insertEnvironmentValues(dsn)
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

func extractComposePatternsFromFile(dockerComposeFile []byte) (out map[string]Pattern, err error) {
	composeServices := map[string]interface{}{}
	{
		// validating prepared config
		err = yaml.Unmarshal(dockerComposeFile, composeServices)
		if err != nil {
			return nil, errors.Wrap(err, "error validating prepared config")
		}

		examplesMap, ok := composeServices[servicesPart]
		if !ok {
			return nil, errors.Wrapf(ErrInvalidComposeFileFormat, "expected to have \"%s\" object", servicesPart)
		}
		composeServices = examplesMap.(map[string]interface{})
	}

	services := make(map[string]Pattern, len(composeServices))

	for serviceName, content := range composeServices {
		cs := Pattern{
			Name: serviceName,
		}

		var bts []byte
		bts, err = yaml.Marshal(content)
		if err != nil {
			return nil, errors.Wrapf(err, "error marshaling service to yaml")
		}

		err = yaml.Unmarshal(bts, &cs.content)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling service "+serviceName+" to struct")
		}

		cs.envs, err = extractEnvsFromComposeFile(bts)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting environment variables")
		}

		services[serviceName] = cs
	}

	return services, nil
}

// extractEnvsFromComposeFile walks thought compose service description
// in order to find environment variable usage such as '${ENV_VAR}'
func extractEnvsFromComposeFile(b []byte) (*env.Container, error) {
	startIdx := bytes.Index(b, []byte{36, 123}) // search for "${"
	var endIdx int

	out := &env.Container{}

	for {
		if startIdx == -1 {
			break
		}

		// Looks like a black magic BUT. NEED this for the 2nd (and others) loop.
		// After next "${" was found, move pointer off last "${" location
		// TODO maybe rewrite it to something more meaningful and easy
		startIdx += endIdx + 1

		endIdx = startIdx + bytes.IndexByte(b[startIdx:], 125) // "}"
		if endIdx == -1 {
			return nil, errors.Wrapf(ErrInvalidComposeEnvFormat, string(b[startIdx:nums.Min(len(b)-1, 10)]))
		}

		// validate if after env variable goes ":" -> extract value
		var val []byte
		if b[endIdx+1] == 58 {
			val = b[endIdx+2:][:bytes.IndexByte(b[endIdx+1:], byte('\n'))]
			if val[len(val)-1] == '\n' {
				val = val[:len(val)-1]
			}
		}

		out.Append(string(b[startIdx+1:endIdx]), string(val))

		startIdx = bytes.Index(b[endIdx:], []byte{36, 123}) // "${"

	}

	return out, nil
}

func removeEnvironmentBrackets(in string) string {
	if in[:2] == "${" && in[len(in)-1] == '}' {
		return in[2 : len(in)-1]
	}

	return in
}
func AddEnvironmentBrackets(in string) string {
	return "${" + in + "}"
}
