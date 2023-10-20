package compose

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/Red-Sock/trace-errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/plugins/environment/project/compose/env"
	"github.com/Red-Sock/rscli/plugins/environment/project/envpatterns"

	"github.com/Red-Sock/rscli/internal/utils/copier"
	"github.com/Red-Sock/rscli/internal/utils/nums"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

const (
	servicesPart = "services"
)

var (
	ErrInvalidComposeFileFormat = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposeEnvFormat  = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

type PatternManager struct {
	Patterns map[string]Pattern
}

type Pattern struct {
	Name                string
	ResourceType        string
	ContainerDefinition ContainerSettings
	Envs                *env.Container
}

func ReadComposePatternsFromFile(pth string) (out *PatternManager, err error) {
	out = &PatternManager{}
	// Basic compose examples: rscli built-in
	out.Patterns, err = extractComposePatternsFromFile(envpatterns.BuildInComposeExamples.Content)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting composePatterns from prepared file")
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

func (p *Pattern) GetName() string {
	return p.Name
}

func (p *Pattern) GetType() string {
	return p.ResourceType
}

func (p *Pattern) GetEnvs() *env.Container {
	return p.Envs
}

func (p *Pattern) GetCompose() ContainerSettings {
	return p.ContainerDefinition
}

func (p *Pattern) RenameVariable(oldName, newName string) {
	p.Envs.Rename(oldName, newName)

	for k, v := range p.ContainerDefinition.Environment {
		if strings.Contains(v, oldName) {
			p.ContainerDefinition.Environment[k] = AddEnvironmentBrackets(newName)
			break
		}
	}

	for portIdx := range p.ContainerDefinition.Ports {
		p.ContainerDefinition.Ports[portIdx] = strings.ReplaceAll(p.ContainerDefinition.Ports[portIdx], oldName, newName)
	}
}

func (c *PatternManager) GetServiceDependencies(cfg *config.Config) ([]Pattern, error) {
	resource, err := cfg.GetDataSourceOptions()
	if err != nil {
		return nil, errors.Wrap(err, "error obtaining data source options")
	}

	clients := make([]Pattern, 0, len(resource))

	for _, resourceDependency := range resource {
		originalPattern, ok := c.Patterns[string(resourceDependency.GetType())]
		if !ok {
			// TODO handle unknown data sources
			continue
		}
		var pattern Pattern
		err := copier.Copy(&originalPattern, &pattern)
		if err != nil {
			return nil, errors.Wrap(err, "error coping pattern")
		}

		if resourceDependency.GetName() != "" {
			pattern.Name += "_" + resourceDependency.GetName()
		}

		// if value for environment variable is defined in source config
		// e.g. section postgres_* in dev.yaml
		// data_sources:
		// 		postgres:
		//			pwd: 123
		//
		// will set environment variable for POSTGRES_PASSWORD
		// by setting variable PROJ_NAME_CAPS_RESOURCE_NAME_CAPS_PWD in .env file to 123
		for name, val := range resourceDependency.GetEnv() {
			if envVariable, ok := pattern.ContainerDefinition.Environment[name]; ok {
				pattern.Envs.AppendRaw(removeEnvironmentBrackets(envVariable), val)
			}
		}

		clients = append(clients, pattern)
	}

	return clients, nil

}

func AddEnvironmentBrackets(in string) string {
	return "${" + in + "}"
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
			Name:         serviceName,
			ResourceType: serviceName,
		}

		var bts []byte
		bts, err = yaml.Marshal(content)
		if err != nil {
			return nil, errors.Wrapf(err, "error marshaling service to yaml")
		}

		err = yaml.Unmarshal(bts, &cs.ContainerDefinition)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling service "+serviceName+" to struct")
		}

		cs.Envs, err = extractEnvsFromComposeFile(bts)
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

		out.AppendRaw(string(b[startIdx+1:endIdx]), string(val))

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
